package wallet

import (
	"context"
	"errors"

	"github.com/3bd-dev/wallet-service/internal/dto/request"
	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/internal/payment"
	"github.com/3bd-dev/wallet-service/pkg/errs"
	"github.com/3bd-dev/wallet-service/pkg/logger"
	"github.com/3bd-dev/wallet-service/pkg/queue"
	"github.com/google/uuid"
)

type Service struct {
	log             *logger.Logger
	walletRepo      IWalletRepo
	transactionRepo ITransactionRepo
	paymentHandler  IPaymentHandler
	cbformat        string
	tranQueue       *queue.Queue[QueueItem]
}

func NewService(log *logger.Logger, walletRepo IWalletRepo, transactionRepo ITransactionRepo, paymenth IPaymentHandler, cbformat string) *Service {
	return &Service{
		log:             log,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		paymentHandler:  paymenth,
		cbformat:        cbformat,
		tranQueue:       queue.NewQueue[QueueItem](),
	}
}

// Deposit creates a new deposit transaction for the given wallet.
func (s *Service) Deposit(ctx context.Context, walletID uuid.UUID, req request.Deposit) (*models.Transaction, error) {
	if err := errs.Check(req); err != nil {
		return nil, err
	}

	paymentMethod, err := s.paymentHandler.VerifyMethod(req.Payment.Gateway, models.TransactionTypeDeposit, req.Payment.Method, req.Payment.MethodDetails)
	if err != nil {
		return nil, err
	}

	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	transaction := &models.Transaction{
		ID:                   uuid.New(),
		WalletID:             wallet.ID,
		Amount:               req.Amount,
		Status:               models.TransactionStatusCreated,
		Type:                 models.TransactionTypeDeposit,
		PaymentGateway:       req.Payment.Gateway,
		PaymentMethodDetails: paymentMethod.MaskRaw(),
		PaymentMethod:        req.Payment.Method,
	}

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, errs.New(errs.Internal, err)
	}

	s.enqueueTransaction(transaction.ID, paymentMethod.GetRaw())

	return transaction, nil
}

// Withdraw creates a new withdrawal transaction for the given wallet.
func (s *Service) Withdraw(ctx context.Context, walletID uuid.UUID, req request.Withdraw) (*models.Transaction, error) {
	if err := errs.Check(req); err != nil {
		return nil, err
	}

	paymentMethod, err := s.paymentHandler.VerifyMethod(req.Payment.Gateway, models.TransactionTypeWithdrawal, req.Payment.Method, req.Payment.MethodDetails)
	if err != nil {
		return nil, err
	}

	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	transaction := &models.Transaction{
		ID:                   uuid.New(),
		WalletID:             wallet.ID,
		Amount:               req.Amount,
		Status:               models.TransactionStatusCreated,
		Type:                 models.TransactionTypeWithdrawal,
		PaymentGateway:       req.Payment.Gateway,
		PaymentMethodDetails: paymentMethod.MaskRaw(),
		PaymentMethod:        req.Payment.Method,
	}

	s.enqueueTransaction(transaction.ID, paymentMethod.GetRaw())

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, errs.New(errs.Internal, err)
	}

	return transaction, nil
}

// ProcessCallback processes the callback from the payment gateway.
func (s *Service) ProcessCallback(ctx context.Context, walletID, tranID uuid.UUID, body []byte) error {
	transaction, err := s.transactionRepo.GetByIDAndWalletID(ctx, tranID, walletID)
	if err != nil {
		return err
	}

	if transaction.Status != models.TransactionStatusPending || transaction.WalletID != walletID {
		return errs.New(errs.InvalidArgument, errors.New("invalid request"))
	}

	var tranRefID string
	if transaction.ReferenceID != nil {
		tranRefID = *transaction.ReferenceID
	}

	res, err := s.paymentHandler.VerifyCallback(ctx, transaction.PaymentGateway, tranRefID, body)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	switch res.Status {
	case payment.PaymentStatusSuccess:
		transaction.Status = models.TransactionStatusCompleted
	case payment.PaymentStatusFailed:
		transaction.Status = models.TransactionStatusFailed
	case payment.PaymentStatusPending:
		return nil
	case payment.PaymentStatusUnknown:
		return errs.New(errs.InvalidArgument, errors.New("unknown payment status"))
	}

	err = s.transactionRepo.Update(ctx, transaction)
	if err != nil {
		return errs.New(errs.Internal, err)
	}
	return nil
}

// GetTransaction retrieves a transaction by its ID and wallet ID.
func (s *Service) GetTransaction(ctx context.Context, id, walletID uuid.UUID) (*models.Transaction, error) {
	tran, err := s.transactionRepo.GetByIDAndWalletID(ctx, id, walletID)
	if err != nil {
		return nil, err
	}
	return tran, nil
}

// Create creates a new wallet.
func (s *Service) Create(ctx context.Context, req request.CreateWallet) (*models.Wallet, error) {
	if err := errs.Check(req); err != nil {
		return nil, err
	}

	wallet := &models.Wallet{
		ID: uuid.New(),
	}

	err := s.walletRepo.Create(ctx, wallet)
	if err != nil {
		return nil, errs.New(errs.Internal, err)
	}
	return wallet, nil
}

// ListTransactions retrieves all transactions for the given wallet.
func (s *Service) ListTransactions(ctx context.Context, walletID uuid.UUID) ([]models.Transaction, error) {
	tran, err := s.transactionRepo.GetByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return tran, nil
}

// List retrieves all wallets.
func (s *Service) List(ctx context.Context) ([]models.Wallet, error) {
	wallets, err := s.walletRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}
