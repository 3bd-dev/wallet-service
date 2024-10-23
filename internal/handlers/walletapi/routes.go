package walletapi

import (
	"net/http"

	"github.com/3bd-dev/wallet-service/internal/services/wallet"
	"github.com/gorilla/mux"
)

type Config struct {
	Service *wallet.Service
}

// Routes adds specific routes for this group.
func Routes(router *mux.Router, cfg Config) {
	api := newapi(cfg.Service)
	wallets := router.PathPrefix("/api/v1/wallets").Subrouter()
	wallets.HandleFunc("", api.create).Methods(http.MethodPost)
	wallets.HandleFunc("", api.list).Methods(http.MethodGet)
	wallets.HandleFunc("/{id}/deposit", api.deposit).Methods(http.MethodPost)
	wallets.HandleFunc("/{id}/withdraw", api.withdraw).Methods(http.MethodPost)
	wallets.HandleFunc("/{id}/transactions", api.getTransactions).Methods(http.MethodGet)
	wallets.HandleFunc("/{id}/transactions/{transactionID}", api.getTransaction).Methods(http.MethodGet)
	wallets.HandleFunc("/{id}/transactions/{transactionID}/callback", api.callback).Methods(http.MethodPost)
}
