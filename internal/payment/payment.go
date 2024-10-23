package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/pkg/errs"
)

// PaymentGateway defines the common interface for all payment gateways
type PaymentGateway interface {
	Deposit(ctx context.Context, req *Request) (*Response, error)
	Withdraw(ctx context.Context, req *Request) (*Response, error)
	VerifyCallback(ctx context.Context, refID string, data []byte) (*Response, error)
	VerifyMethod(typ models.TransactionType, method models.PaymentMethod) error
}

type Payment struct {
	gateways map[models.PaymentGateway]PaymentGateway
}

// NewPayment creates a new Payment struct with the provided gateways
func New(gateways map[models.PaymentGateway]PaymentGateway) *Payment {
	return &Payment{
		gateways: gateways,
	}
}

// Deposit sends a deposit request to the appropriate gateway
func (p *Payment) Deposit(ctx context.Context, gateway models.PaymentGateway, req *Request) (*Response, error) {
	if err := p.validateGateway(gateway); err != nil {
		return nil, err
	}

	res, err := p.gateways[gateway].Deposit(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to deposit: %w", err)
	}

	return res, nil
}

// Withdraw sends a withdrawal request to the appropriate gateway
func (p *Payment) Withdraw(ctx context.Context, gateway models.PaymentGateway, req *Request) (*Response, error) {
	if err := p.validateGateway(gateway); err != nil {
		return nil, err
	}

	res, err := p.gateways[gateway].Withdraw(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw: %w", err)
	}

	return res, nil
}

// VerifyCallback chooses the appropriate gateway to handle the callback
func (p *Payment) VerifyCallback(ctx context.Context, gatewayName models.PaymentGateway, refID string, data []byte) (*Response, error) {
	if err := p.validateGateway(gatewayName); err != nil {
		return nil, err
	}

	res, err := p.gateways[gatewayName].VerifyCallback(ctx, refID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to verify callback: %w", err)
	}

	return res, nil
}

// VerifyMethod verifies the payment method details
func (p *Payment) VerifyMethod(gateway models.PaymentGateway, typ models.TransactionType, method models.PaymentMethod, paymMethDet json.RawMessage) (PaymentMethodDetails, error) {
	if err := p.validateGateway(gateway); err != nil {
		return nil, err
	}
	err := p.gateways[gateway].VerifyMethod(typ, method)
	if err != nil {
		return nil, errs.New(errs.InvalidArgument, fmt.Errorf("failed to verify method: %w", err))
	}

	paymMeth, err := p.parsePaymentMethodDetails(method, paymMethDet)
	if err != nil {
		return nil, err
	}

	if err := paymMeth.validate(); err != nil {
		return nil, err
	}
	return paymMeth, nil
}

// parsePaymentMethodDetails decouples the parsing logic to make the addition of new payment methods easier.
func (p *Payment) parsePaymentMethodDetails(method models.PaymentMethod, data json.RawMessage) (PaymentMethodDetails, error) {
	switch method {
	case models.PaymentMethodCreditCard:
		var details PaymentMethodCreditCardDetails
		if err := json.Unmarshal(data, &details); err != nil {
			return nil, errs.New(errs.InvalidArgument, fmt.Errorf("failed to unmarshal credit card details: %w", err))
		}
		return &details, nil
	case models.PaymentMethodBankTransfer:
		var details PaymentMethodBankDetails
		if err := json.Unmarshal(data, &details); err != nil {
			return nil, errs.New(errs.InvalidArgument, fmt.Errorf("failed to unmarshal credit card details: %w", err))
		}
		return &details, nil
	default:
		return nil, errs.New(errs.InvalidArgument, errors.New("unsupported payment method"))
	}
}

func (p *Payment) validateGateway(gateway models.PaymentGateway) error {
	if _, ok := p.gateways[gateway]; !ok {
		return fmt.Errorf("unsupported gateway: %s", gateway)
	}
	return nil
}
