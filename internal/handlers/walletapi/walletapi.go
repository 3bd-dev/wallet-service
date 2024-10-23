package walletapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/3bd-dev/wallet-service/internal/dto/request"
	"github.com/3bd-dev/wallet-service/internal/services/wallet"
	"github.com/3bd-dev/wallet-service/pkg/errs"
	"github.com/3bd-dev/wallet-service/pkg/web"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type api struct {
	service *wallet.Service
}

func newapi(svc *wallet.Service) *api {
	return &api{service: svc}
}

// deposit adds a deposit transaction to the wallet.
func (a *api) deposit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("invalid ID: %w", err)))
		return
	}

	var req request.Deposit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("failed to decode request body: %w", err)))
		return
	}

	tran, er := a.service.Deposit(r.Context(), id, req)
	if er != nil {
		web.RenderErr(w, er)
		return
	}

	web.RenderOk(w, tran)
}

// withdraw adds a withdraw transaction to the wallet.
func (a *api) withdraw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("invalid ID: %w", err)))
		return
	}

	var req request.Withdraw
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("failed to decode request body: %w", err)))
		return
	}

	tran, er := a.service.Withdraw(r.Context(), id, req)
	if er != nil {
		web.RenderErr(w, er)
		return
	}

	web.RenderOk(w, tran)
}

// callback processes the callback from the payment gateway.
func (a *api) callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("invalid ID: %w", err)))
		return
	}

	tranID, err := uuid.Parse(vars["transactionID"])
	if err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("invalid transaction ID: %w", err)))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("failed to read request body: %w", err)))
		return
	}

	if err := a.service.ProcessCallback(r.Context(), id, tranID, body); err != nil {
		web.RenderErr(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getTransaction returns the transaction by ID.
func (a *api) getTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("invalid ID: %w", err)))
		return
	}

	tranID, err := uuid.Parse(vars["transactionID"])
	if err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("invalid transaction ID: %w", err)))
		return
	}

	tran, er := a.service.GetTransaction(r.Context(), tranID, id)
	if er != nil {
		web.RenderErr(w, er)
		return
	}

	web.RenderOk(w, tran)
}

// create creates a new wallet.
func (a *api) create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateWallet
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("failed to decode request body: %w", err)))
		return
	}

	wallet, err := a.service.Create(r.Context(), req)
	if err != nil {
		web.RenderErr(w, err)
		return
	}

	web.RenderOk(w, wallet)
}

// getTransactions returns the transactions for the wallet.
func (a *api) getTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		web.RenderErr(w, errs.New(errs.InvalidArgument, fmt.Errorf("invalid ID: %w", err)))
		return
	}
	tran, er := a.service.ListTransactions(r.Context(), id)
	if er != nil {
		web.RenderErr(w, er)
		return
	}

	web.RenderOk(w, tran)
}

func (a *api) list(w http.ResponseWriter, r *http.Request) {
	wallets, err := a.service.List(r.Context())
	if err != nil {
		web.RenderErr(w, err)
		return
	}

	web.RenderOk(w, wallets)
}
