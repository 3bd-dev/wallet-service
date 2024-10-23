package request

import (
	"encoding/json"

	"github.com/3bd-dev/wallet-service/internal/models"
)

type Payment struct {
	Gateway       models.PaymentGateway `json:"gateway" validate:"required"`
	Method        models.PaymentMethod  `json:"method" validate:"required"`
	MethodDetails json.RawMessage       `json:"method_details" validate:"required,json"`
}

type Deposit struct {
	Amount  float64 `json:"amount" validate:"required,gt=0"`
	Payment Payment `json:"payment"`
}

type Withdraw struct {
	Amount  float64 `json:"amount" validate:"required,gt=0"`
	Payment Payment `json:"payment"`
}
