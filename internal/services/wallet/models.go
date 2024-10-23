package wallet

import (
	"encoding/json"

	"github.com/google/uuid"
)

type QueueItem struct {
	ID             uuid.UUID       `json:"id"`
	PaymentDetails json.RawMessage `json:"payment_details"`
}
