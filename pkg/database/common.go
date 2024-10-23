package database

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

const (
	ContextKeyDBTx ContextKey = "tx"
)
