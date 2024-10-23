package gatewayb

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"github.com/3bd-dev/wallet-service/config"
	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/internal/payment"
	"github.com/3bd-dev/wallet-service/pkg/errs"
	"github.com/3bd-dev/wallet-service/pkg/rest"
	"github.com/sony/gobreaker/v2"
)

var supportedMethod = map[models.TransactionType][]models.PaymentMethod{
	models.TransactionTypeDeposit: {
		models.PaymentMethodCreditCard,
	},
	models.TransactionTypeWithdrawal: {
		models.PaymentMethodBankTransfer,
	},
}

// GatewayB is a concrete implementation of the PaymentGateway To Gateway B
type GatewayB struct {
	client *rest.Client
	retier rest.Retrier
	cb     *gobreaker.CircuitBreaker[[]byte]
}

// NewGatewayB creates a new instance of Gateway B
func New(cfg config.PaymentGatewayB) payment.PaymentGateway {
	cbSettings := gobreaker.Settings{
		Name:        "GatewayB",
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > cfg.CBMaxConsecutiveFailures
		},
	}
	return &GatewayB{
		client: rest.NewClient(cfg.BaseURL),
		retier: rest.NewRetrier(cfg.RetryAttempt, cfg.RetryDelay),
		cb:     gobreaker.NewCircuitBreaker[[]byte](cbSettings),
	}
}

// Deposit sends a deposit request to Gateway B using SOAP/XML
func (g *GatewayB) Deposit(ctx context.Context, req *payment.Request) (*payment.Response, error) {

	body, err := g.cb.Execute(func() ([]byte, error) {
		req := &Request{
			Amount:      req.Amount,
			CallbackURL: req.CallbackURL,
		}

		resp, err := g.retry(ctx, "/deposit", req, &rest.RequestOptions{
			Headers: http.Header{
				"Content-Type": []string{rest.XMLContentType},
			},
		}, g.client.Post)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to initiate deposit: %d - %s", resp.StatusCode, resp.Body)
		}

		return resp.Body, nil
	})

	if err != nil {
		return nil, err
	}

	var result Response
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &payment.Response{ID: result.Body.ReferenceID, Status: toPaymentStatus(result.Body.Status)}, nil
}

// Withdraw sends a withdrawal request to Gateway B using SOAP/XML
func (g *GatewayB) Withdraw(ctx context.Context, req *payment.Request) (*payment.Response, error) {

	body, err := g.cb.Execute(func() ([]byte, error) {
		req := &Request{
			Amount:      req.Amount,
			CallbackURL: req.CallbackURL,
		}

		resp, err := g.retry(ctx, "/withdraw", req, &rest.RequestOptions{
			Headers: http.Header{
				"Content-Type": []string{rest.XMLContentType},
			},
		}, g.client.Post)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to initiate withdraw: %d - %s", resp.StatusCode, resp.Body)
		}

		return resp.Body, nil
	})

	if err != nil {
		return nil, err
	}

	var result Response
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &payment.Response{ID: result.Body.ReferenceID, Status: toPaymentStatus(result.Body.Status)}, nil
}

// VerifyCallback verifies the callback from Gateway B
func (g *GatewayB) VerifyCallback(ctx context.Context, refID string, data []byte) (*payment.Response, error) {
	var res Response
	if err := xml.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	if res.Body.Status == "" {
		return nil, errors.New("failed to process callback")
	}

	if refID != res.Body.ReferenceID {
		return nil, errs.New(errs.InvalidArgument, errors.New("invalid reference ID"))
	}

	return &payment.Response{ID: res.Body.ReferenceID, Status: toPaymentStatus(res.Body.Status)}, nil
}

// VerifyMethod verifies the payment method details
func (g *GatewayB) VerifyMethod(typ models.TransactionType, method models.PaymentMethod) error {
	if methods, ok := supportedMethod[typ]; ok {
		for _, m := range methods {
			if m == method {
				return nil
			}
		}
	}
	return errs.New(errs.InvalidArgument, errors.New("unsupported payment method"))
}

// retry retries the request if it fails
func (g *GatewayB) retry(ctx context.Context, url string, body any, options *rest.RequestOptions, fn func(ctx context.Context, reqURL string, body interface{}, options *rest.RequestOptions) (*rest.Response, error)) (*rest.Response, error) {
	var resp *rest.Response
	err := g.retier.Do(ctx, func(ctx context.Context, _ int) error {
		var err error
		resp, err = fn(ctx, url, body, options)
		if err != nil {
			return err
		}

		if resp == nil || resp.StatusCode >= http.StatusInternalServerError {
			return errors.New("invalid response")
		}

		return err
	})

	return resp, err
}

func toPaymentStatus(status string) payment.PaymentStatus {
	switch status {
	case "success":
		return payment.PaymentStatusSuccess
	case "pending":
		return payment.PaymentStatusPending
	case "failed":
		return payment.PaymentStatusFailed
	default:
		return payment.PaymentStatusUnknown
	}
}
