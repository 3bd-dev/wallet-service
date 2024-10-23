package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// GORM model definition
type Transaction struct {
	ID                   uuid.UUID         `json:"id"`
	WalletID             uuid.UUID         `json:"wallet_id"`
	Amount               float64           `json:"amount"`
	Type                 TransactionType   `json:"type"`
	Status               TransactionStatus `json:"status"`
	PaymentGateway       PaymentGateway    `json:"payment_gateway"`
	PaymentMethod        PaymentMethod     `json:"payment_method"`
	PaymentMethodDetails json.RawMessage   `json:"payment_method_details"`
	ReferenceID          *string           `json:"reference_id"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
	Wallet               *Wallet           `json:"wallet,omitempty" `
}

func (t *Transaction) IsEmpty() bool {
	return t == nil || t.ID == uuid.Nil
}

// -------------
//
// TransactionType represents the type of a transaction
type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusCreated   TransactionStatus = "created"
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// PaymentGateway represents the payment gateway used for a transaction
type PaymentGateway string

const (
	PaymentGatewayA PaymentGateway = "gateway_a"
	PaymentGatewayB PaymentGateway = "gateway_b"
)

// PaymentMethod represents the payment method used for a transaction
type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)
