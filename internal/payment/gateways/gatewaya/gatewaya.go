package gatewaya

import (
	"context"
	"encoding/json"
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

// GatewayA represents the Gateway A payment gateway
type GatewayA struct {
	client *rest.Client
	retier rest.Retrier
	cb     *gobreaker.CircuitBreaker[[]byte]
}

// NewGateway creates a new instance of Gateway A
func New(cfg config.PaymentGatewayA) payment.PaymentGateway {
	cbSettings := gobreaker.Settings{
		Name:        "GatewayA",
		MaxRequests: cfg.CBMaxRequests,
		Interval:    cfg.CBInterval,
		Timeout:     cfg.CBTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > cfg.CBMaxConsecutiveFailures || counts.TotalFailures > cfg.CBMaxTotalFailures
		},
	}
	return &GatewayA{
		client: rest.NewClient(cfg.BaseURL),
		retier: rest.NewRetrier(cfg.RetryAttempt, cfg.RetryDelay),
		cb:     gobreaker.NewCircuitBreaker[[]byte](cbSettings),
	}
}

// Deposit sends a deposit request to Gateway A
func (g *GatewayA) Deposit(ctx context.Context, req *payment.Request) (*payment.Response, error) {
	body, err := g.cb.Execute(func() ([]byte, error) {
		requestBody := Request{
			Amount:      req.Amount,
			CallbackURL: req.CallbackURL,
		}
		resp, err := g.retry(ctx, "/deposit", requestBody, nil, g.client.Post)
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

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &payment.Response{ID: res.ID, Status: toPaymentStatus(res.Status)}, nil
}

// Withdrawal sends a withdrawal request to Gateway A
func (g *GatewayA) Withdraw(ctx context.Context, req *payment.Request) (*payment.Response, error) {
	body, err := g.cb.Execute(func() ([]byte, error) {
		requestBody := Request{
			Amount:      req.Amount,
			CallbackURL: req.CallbackURL,
		}

		resp, err := g.retry(ctx, "/withdrawal", requestBody, nil, g.client.Post)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to initiate withdraw: : %d - %s", resp.StatusCode, string(resp.Body))
		}

		return resp.Body, nil
	})

	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &payment.Response{ID: res.ID, Status: toPaymentStatus(res.Status)}, nil
}

// VerifyCallback processes the callback from Gateway A
func (g *GatewayA) VerifyCallback(ctx context.Context, refID string, data []byte) (*payment.Response, error) {
	var res Response
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	if res.Status == "" {
		return nil, errors.New("failed to process transaction")
	}

	if refID != res.ID {
		return nil, errs.New(errs.InvalidArgument, errors.New("invalid reference ID"))
	}

	return &payment.Response{ID: res.ID, Status: toPaymentStatus(res.Status)}, nil
}

// VerifyMethod verifies the payment method details
func (g *GatewayA) VerifyMethod(typ models.TransactionType, method models.PaymentMethod) error {
	if methods, ok := supportedMethod[typ]; ok {
		for _, m := range methods {
			if m == method {
				return nil
			}
		}
	}
	return errs.New(errs.InvalidArgument, errors.New("unsupported payment method"))
}

// retry sends a request to the gateway and retries if it fails
func (g *GatewayA) retry(ctx context.Context, url string, body any, options *rest.RequestOptions, fn func(ctx context.Context, reqURL string, body interface{}, options *rest.RequestOptions) (*rest.Response, error)) (*rest.Response, error) {
	var resp *rest.Response
	err := g.retier.Do(ctx, func(ctx context.Context, i int) error {
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
