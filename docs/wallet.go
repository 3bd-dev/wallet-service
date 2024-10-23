package docs

import (
	"github.com/3bd-dev/wallet-service/internal/dto/request"
	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/pkg/web"
	"github.com/google/uuid"
)

// swagger:route GET /api/v1/wallets Wallets ListWallets
// List all cities
// responses:
//   200: ListWalletsResponse

// swagger:response ListWalletsResponse
type ListWalletsResponseWrapper struct {
	// in:body
	Body struct {
		web.Response
		Data []models.Wallet `json:"data"`
	}
}

// swagger:route POST /api/v1/wallets Wallets CreateWallet
// Create a new wallet
// responses:
//   201: CreateWalletResponse

// swagger:parameters CreateWallet
type CreateWalletParamsWrapper struct {
	// in:body
	Body struct {
	}
}

// swagger:response CreateWalletResponse
type CreateWalletResponseWrapper struct {
	// in:body
	Body struct {
		web.Response
		Data models.Wallet `json:"data"`
	}
}

// swagger:route GET /api/v1/wallets/{id} Wallets GetWallet
// Get a wallet by id
// responses:
//   200: GetWalletResponse

// swagger:parameters GetWallet
type GetWalletParamsWrapper struct {
	// in:path
	// Required: true
	ID uuid.UUID `json:"id"`
}

// swagger:response GetWalletResponse
type GetWalletResponseWrapper struct {
	// in:body
	Body struct {
		web.Response
		Data models.Wallet `json:"data"`
	}
}

// swagger:route Post /api/v1/wallets/{id}/deposit Wallets MakeDeposit
// Deposit to a wallet by id
// responses:
//   200: DepositResponse

// swagger:parameters MakeDeposit
type DepositParamsWrapper struct {
	// in:path
	// Required: true
	ID uuid.UUID `json:"id"`
	// in:body
	Body struct {
		request.Deposit
	}
}

// swagger:response DepositResponse
type DepositResponseWrapper struct {
	// in:body
	Body struct {
		web.Response
		Data models.Wallet `json:"data"`
	}
}

// swagger:route Post /api/v1/wallets/{id}/withdraw Wallets MakeWithdraw
// Withdraw from a wallet by id
// responses:
//   200: WithdrawResponse

// swagger:parameters MakeWithdraw
type WithdrawParamsWrapper struct {
	// in:path
	// Required: true
	ID uuid.UUID `json:"id"`
	// in:body
	Body struct {
		request.Withdraw
	}
}

// swagger:response WithdrawResponse
type WithdrawResponseWrapper struct {
	// in:body
	Body struct {
		web.Response
		Data models.Wallet `json:"data"`
	}
}

// swagger:route Get /api/v1/wallets/{id}/transactions Transactions ListTransactions
// List transactions of a wallet by id
// responses:
//   200: ListTransactionsResponse

// swagger:parameters ListTransactions
type ListTransactionsParamsWrapper struct {
	// in:path
	// Required: true
	ID uuid.UUID `json:"id"`
}

// swagger:response ListTransactionsResponse
type ListTransactionsResponseWrapper struct {
	// in:body
	Body struct {
		web.Response
		Data []models.Transaction `json:"data"`
	}
}

// swagger:route Get /api/v1/wallets/{id}/transactions/{transaction_id} Transactions GetTransaction
// Get a transaction of a wallet by id
// responses:
//   200: GetTransactionResponse

// swagger:parameters GetTransaction
type GetTransactionParamsWrapper struct {
	// in:path
	// Required: true
	ID uuid.UUID `json:"id"`
	// in:path
	// Required: true
	TransactionID uuid.UUID `json:"transaction_id"`
}

// swagger:response GetTransactionResponse
type GetTransactionResponseWrapper struct {
	// in:body
	Body struct {
		web.Response
		Data models.Transaction `json:"data"`
	}
}
