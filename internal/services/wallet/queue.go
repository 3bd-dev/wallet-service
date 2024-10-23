package wallet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/internal/payment"
	"github.com/google/uuid"
)

// processTransaction handles the processing of a single transaction, with retry logic and error handling.
func (s *Service) processTransaction(ctx context.Context, item QueueItem) {
	s.log.Debug(ctx, "Processing transaction", "transaction_id", item.ID)
	var err error
	var res *payment.Response
	var tran *models.Transaction

	defer func() {
		if err != nil {
			s.log.Error(ctx, "Failed to process transaction", "transaction_id", item.ID, "error", err)
			if tran != nil {
				tran.Status = models.TransactionStatusFailed
			}
		}

		if tran.Status == models.TransactionStatusFailed {
			if updateErr := s.transactionRepo.Update(ctx, tran); updateErr != nil {
				s.log.Error(ctx, "Failed to update transaction status", "transaction_id", item.ID, "error", updateErr)
			}
		}
	}()

	tran, err = s.transactionRepo.GetByID(ctx, item.ID)
	if err != nil {
		return
	}

	paymentReq := &payment.Request{
		ID:                   tran.ID.String(),
		Amount:               tran.Amount,
		CallbackURL:          fmt.Sprintf(s.cbformat, tran.WalletID, tran.ID),
		PaymentMethodDetails: item.PaymentDetails,
	}

	// Process the transaction based on its type
	switch tran.Type {
	case models.TransactionTypeDeposit:
		res, err = s.paymentHandler.Deposit(ctx, tran.PaymentGateway, paymentReq)
	case models.TransactionTypeWithdrawal:
		res, err = s.paymentHandler.Withdraw(ctx, tran.PaymentGateway, paymentReq)
	default:
		err = fmt.Errorf("unsupported transaction type: %s", tran.Type)
	}

	if err != nil {
		return
	}

	tran.ReferenceID = &res.ID
	tran.Status = models.TransactionStatusPending

	if err = s.transactionRepo.Update(ctx, tran); err != nil {
		return
	}

	s.log.Info(ctx, "Transaction processed successfully", "transaction_id", tran.ID)
}

// enqueueTransaction adds a transaction to the queue for processing.
func (s *Service) enqueueTransaction(tranID uuid.UUID, paymentDetails json.RawMessage) {
	s.tranQueue.Enqueue(QueueItem{
		ID:             tranID,
		PaymentDetails: paymentDetails,
	})
}

// Start starts the transaction queue worker.
func (s *Service) Start(ctx context.Context) {
	s.tranQueue.StartWorker(ctx, s.processTransaction)
}
