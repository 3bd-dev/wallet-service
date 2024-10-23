package payment

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/pkg/errs"
	"github.com/3bd-dev/wallet-service/pkg/unitest"
)

type mockGateway struct {
	depositFunc        func(ctx context.Context, req *Request) (*Response, error)
	withdrawFunc       func(ctx context.Context, req *Request) (*Response, error)
	verifyCallbackFunc func(ctx context.Context, refID string, data []byte) (*Response, error)
	verifyMethodFunc   func(typ models.TransactionType, method models.PaymentMethod) error
}

func (m *mockGateway) Deposit(ctx context.Context, req *Request) (*Response, error) {
	return m.depositFunc(ctx, req)
}

func (m *mockGateway) Withdraw(ctx context.Context, req *Request) (*Response, error) {
	return m.withdrawFunc(ctx, req)
}

func (m *mockGateway) VerifyCallback(ctx context.Context, refID string, data []byte) (*Response, error) {
	return m.verifyCallbackFunc(ctx, refID, data)
}

func (m *mockGateway) VerifyMethod(typ models.TransactionType, method models.PaymentMethod) error {
	return m.verifyMethodFunc(typ, method)
}
func Test_Payment(t *testing.T) {
	t.Parallel()

	unitest.Run(t, deposit(), "deposit")
	unitest.Run(t, withdrawal(), "withdrawal")
	unitest.Run(t, VerifyCallback(), "verifyCallback")
	unitest.Run(t, verifyMethodBankTransfer(), "verifyMethodBankTransfer")
	unitest.Run(t, verifyMethodCreditCard(), "verifyMethodCreditCard")
}

func deposit() []unitest.Table {
	mockGateway := &mockGateway{
		depositFunc: func(ctx context.Context, req *Request) (*Response, error) {
			return &Response{Status: "pending"}, nil
		},
	}

	payment := New(map[models.PaymentGateway]PaymentGateway{
		models.PaymentGateway("mock"): mockGateway,
	})

	tests := []unitest.Table{
		{
			Name: "Successful Deposit",
			ExpResp: &Response{
				Status: "pending",
			},
			ExcFunc: func(ctx context.Context) any {
				resp, _ := payment.Deposit(ctx, models.PaymentGateway("mock"), &Request{})
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*Response)
				expResp := exp.(*Response)
				if gotResp.Status != expResp.Status {
					return "status mismatch"
				}
				return ""
			},
		},
		{
			Name:    "Unsupported Gateway",
			ExpResp: errors.New("unsupported gateway: invalid"),
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.Deposit(ctx, models.PaymentGateway("invalid"), &Request{})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(error)
				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return "error message mismatch"
				}
				return ""
			},
		},
	}

	return tests
}

func withdrawal() []unitest.Table {
	mockGateway := &mockGateway{
		withdrawFunc: func(ctx context.Context, req *Request) (*Response, error) {
			return &Response{Status: "success"}, nil
		},
	}

	payment := New(map[models.PaymentGateway]PaymentGateway{
		models.PaymentGateway("mock"): mockGateway,
	})

	tests := []unitest.Table{
		{
			Name: "Successful Withdraw",
			ExpResp: &Response{
				Status: "success",
			},
			ExcFunc: func(ctx context.Context) any {
				resp, _ := payment.Withdraw(ctx, models.PaymentGateway("mock"), &Request{})
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*Response)
				expResp := exp.(*Response)
				if gotResp.Status != expResp.Status {
					return "status mismatch"
				}
				return ""
			},
		},
		{
			Name:    "Unsupported Gateway",
			ExpResp: errors.New("unsupported gateway: invalid"),
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.Withdraw(ctx, models.PaymentGateway("invalid"), &Request{})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(error)
				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return "error message mismatch"
				}
				return ""
			},
		},
	}

	return tests
}

func VerifyCallback() []unitest.Table {
	mockGateway := &mockGateway{
		verifyCallbackFunc: func(ctx context.Context, refID string, data []byte) (*Response, error) {
			return &Response{Status: "success"}, nil
		},
	}

	payment := New(map[models.PaymentGateway]PaymentGateway{
		models.PaymentGateway("mock"): mockGateway,
	})

	tests := []unitest.Table{
		{
			Name: "Successful VerifyCallback",
			ExpResp: &Response{
				Status: "success",
			},
			ExcFunc: func(ctx context.Context) any {
				resp, _ := payment.VerifyCallback(ctx, models.PaymentGateway("mock"), "ref123", []byte{})
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*Response)
				expResp := exp.(*Response)
				if gotResp.Status != expResp.Status {
					return "status mismatch"
				}
				return ""
			},
		},
		{
			Name:    "Unsupported Gateway",
			ExpResp: errors.New("unsupported gateway: invalid"),
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.VerifyCallback(ctx, models.PaymentGateway("invalid"), "ref123", []byte{})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(error)
				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return "error message mismatch"
				}
				return ""
			},
		},
	}

	return tests
}

func verifyMethodBankTransfer() []unitest.Table {
	mockGateway := &mockGateway{
		verifyMethodFunc: func(typ models.TransactionType, method models.PaymentMethod) error {
			if typ == models.TransactionType("valid") && method == models.PaymentMethodBankTransfer {
				return nil
			}
			return fmt.Errorf("invalid transaction type or payment method")
		},
	}

	payment := New(map[models.PaymentGateway]PaymentGateway{
		models.PaymentGateway("mock"): mockGateway,
	})

	validCreditCard := []byte(`{"number": "4111111111111111", "expiry": "12/26", "cvv": "123"}`)
	validBandTransfer := []byte(`{"account_number": "1234567890", "bank_code": "BOFAUS3NXXX","bank_code_type": "SWIFT"}`)

	tests := []unitest.Table{
		{
			Name:    "Invalid Payment Method Details",
			ExpResp: errs.InvalidArgument,
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.VerifyMethod(models.PaymentGateway("mock"), models.TransactionType("valid"), models.PaymentMethodBankTransfer, []byte{})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				if got != nil {
					gotErr := got.(*errs.Error)
					expRes := exp.(errs.ErrCode)
					if gotErr.Code != expRes {
						return fmt.Sprintf("expected error code %v, got %v", errs.InvalidArgument, gotErr.Code)
					}
				}
				return ""
			},
		}, {
			Name:    "InValid Payment Method Credit Card",
			ExpResp: errs.InvalidArgument,
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.VerifyMethod(models.PaymentGateway("mock"), models.TransactionType("valid"), models.PaymentMethodCreditCard, validCreditCard)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expRes := exp.(errs.ErrCode)
				if gotErr.Code != expRes {
					return fmt.Sprintf("expected error code %v, got %v", errs.InvalidArgument, got)
				}
				return ""
			},
		},
		{
			Name:    "Valid Payment Method Details Bank Transfer",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.VerifyMethod(models.PaymentGateway("mock"), models.TransactionType("valid"), models.PaymentMethodBankTransfer, validBandTransfer)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				if got != nil {
					return fmt.Sprintf("expected nil, got %v", got)
				}
				return ""
			},
		},
		{
			Name:    "Valid Payment Method Details Credit Card",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.VerifyMethod(models.PaymentGateway("mock"), models.TransactionType("valid"), models.PaymentMethodBankTransfer, validBandTransfer)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				if got != nil {
					return fmt.Sprintf("expected nil, got %v", got)
				}
				return ""
			},
		},
	}

	return tests
}

func verifyMethodCreditCard() []unitest.Table {
	mockGateway := &mockGateway{
		verifyMethodFunc: func(typ models.TransactionType, method models.PaymentMethod) error {
			if typ == models.TransactionType("valid") && method == models.PaymentMethodCreditCard {
				return nil
			}
			return fmt.Errorf("invalid transaction type or payment method")
		},
	}

	payment := New(map[models.PaymentGateway]PaymentGateway{
		models.PaymentGateway("mock"): mockGateway,
	})

	validCreditCard := []byte(`{"number": "4111111111111111", "expiry": "12/26", "cvv": "123"}`)
	validBandTransfer := []byte(`{"account_number": "1234567890", "bank_code": "BOFAUS3NXXX","bank_code_type": "SWIFT"}`)

	tests := []unitest.Table{
		{
			Name:    "Invalid Payment Method Details",
			ExpResp: errs.InvalidArgument,
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.VerifyMethod(models.PaymentGateway("mock"), models.TransactionType("valid"), models.PaymentMethodCreditCard, []byte{})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				if got != nil {
					gotErr := got.(*errs.Error)
					expRes := exp.(errs.ErrCode)
					if gotErr.Code != expRes {
						return fmt.Sprintf("expected error code %v, got %v", errs.InvalidArgument, gotErr.Code)
					}
				}
				return ""
			},
		}, {
			Name:    "InValid Payment Method Bank Transfer",
			ExpResp: errs.InvalidArgument,
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.VerifyMethod(models.PaymentGateway("mock"), models.TransactionType("valid"), models.PaymentMethodBankTransfer, validBandTransfer)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expRes := exp.(errs.ErrCode)
				if gotErr.Code != expRes {
					return fmt.Sprintf("expected error code %v, got %v", errs.InvalidArgument, got)
				}
				return ""
			},
		},
		{
			Name:    "Valid Payment Method Details Credit Card",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				_, err := payment.VerifyMethod(models.PaymentGateway("mock"), models.TransactionType("valid"), models.PaymentMethodCreditCard, validCreditCard)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				if got != nil {
					return fmt.Sprintf("expected nil, got %v", got)
				}
				return ""
			},
		},
	}

	return tests
}
