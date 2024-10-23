package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID           uuid.UUID     `json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

func (t *Wallet) IsEmpty() bool {
	return t == nil || t.ID == uuid.Nil
}
