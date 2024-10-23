package wallet

import (
	"context"
	"encoding/json"

	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/internal/payment"
	"github.com/google/uuid"
)

// TransactionRepo defines the interface for transaction repository.
type ITransactionRepo interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetByIDAndWalletID(ctx context.Context, id, walletID uuid.UUID) (*models.Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	Update(ctx context.Context, wallet *models.Transaction) error
	GetByWalletID(ctx context.Context, walletID uuid.UUID) ([]models.Transaction, error)
}

type IWalletRepo interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	List(ctx context.Context) ([]models.Wallet, error)
}

type IPaymentHandler interface {
	Deposit(ctx context.Context, gateway models.PaymentGateway, req *payment.Request) (*payment.Response, error)
	Withdraw(ctx context.Context, gateway models.PaymentGateway, req *payment.Request) (*payment.Response, error)
	VerifyCallback(ctx context.Context, gatewayName models.PaymentGateway, tranID string, data []byte) (*payment.Response, error)
	VerifyMethod(gateway models.PaymentGateway, typ models.TransactionType, method models.PaymentMethod, paymMethDet json.RawMessage) (payment.PaymentMethodDetails, error)
}
