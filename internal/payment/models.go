package payment

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/3bd-dev/wallet-service/pkg/errs"
	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusUnknown PaymentStatus = "unknown"
)

type Request struct {
	ID                   string          `json:"id"`
	Amount               float64         `json:"amount"`
	CallbackURL          string          `json:"callback_url"`
	PaymentMethodDetails json.RawMessage `json:"payment_details"`
}

type Response struct {
	Status PaymentStatus `json:"status"`
	ID     string        `json:"id"`
}

type PaymentMethodCreditCardDetails struct {
	Number string `json:"number" validate:"required,credit_card"`
	Expiry string `json:"expiry" validate:"required"`
	CVV    string `json:"cvv" validate:"required,numeric,len=3"`
}

type PaymentMethodDetails interface {
	GetRaw() json.RawMessage
	MaskRaw() json.RawMessage
	validate() error
}

func (p *PaymentMethodCreditCardDetails) GetRaw() json.RawMessage {
	// convert to json.RawMessage to avoid infinite recursion
	b, _ := json.Marshal(p)
	return b
}

func (p *PaymentMethodCreditCardDetails) MaskRaw() json.RawMessage {
	return (&PaymentMethodCreditCardDetails{
		Number: fmt.Sprintf("**** **** **** %s", p.Number[len(p.Number)-4:]),
		Expiry: p.Expiry,
		CVV:    "***",
	}).GetRaw()
}

func (p *PaymentMethodCreditCardDetails) validate() error {
	err := errs.Check(p)
	if err != nil {
		return err
	}

	expirationParts := strings.Split(p.Expiry, "/")
	if len(expirationParts) != 2 {
		return fmt.Errorf("invalid expiration date format, expected MM/YY")
	}

	month, err := strconv.Atoi(expirationParts[0])
	if err != nil || month < 1 || month > 12 {
		return fmt.Errorf("invalid expiration month: %s", expirationParts[0])
	}

	year, err := strconv.Atoi(expirationParts[1])
	if err != nil {
		return fmt.Errorf("invalid expiration year: %s", expirationParts[1])
	}

	expirationDate := time.Date(2000+year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	currentTime := time.Now()
	minValidExpiration := currentTime.AddDate(0, 6, 0) // 6 months from now
	if expirationDate.Before(minValidExpiration) {
		return fmt.Errorf("expiration date must be at least 6 months from now")
	}

	return nil
}

type PaymentMethodBankDetails struct {
	AccountNumber string `json:"account_number" validate:"required,numeric,min=10,max=34"` // Account number: 10-34 digits
	BankCode      string `json:"bank_code" validate:"required,alphanum,min=6,max=34"`      // Bank code: varies depending on type (SWIFT, IBAN, etc.)
	BankCodeType  string `json:"bank_code_type" validate:"required,oneof=SWIFT IBAN ROUTING SORTCODE IFSC CLABE"`
}

func (p *PaymentMethodBankDetails) validate() error {
	err := errs.Check(p)
	if err != nil {
		return err
	}

	return nil
}
func (p *PaymentMethodBankDetails) GetRaw() json.RawMessage {
	b, _ := json.Marshal(p)
	return b
}

func (p *PaymentMethodBankDetails) MaskRaw() json.RawMessage {
	return (&PaymentMethodBankDetails{
		AccountNumber: fmt.Sprintf("****%s", p.AccountNumber[len(p.AccountNumber)-4:]),
		BankCode:      fmt.Sprintf("****%s", p.BankCode[len(p.BankCode)-4:]),
		BankCodeType:  p.BankCodeType,
	}).GetRaw()
}

type PaymentRequest struct {
	ID                   uuid.UUID
	PaymentMethodDetails PaymentMethodDetails
	Amount               float64
	CallbackURL          string
}
