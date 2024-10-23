package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/pkg/database"
	"github.com/3bd-dev/wallet-service/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepo struct {
	db database.IDatabase
}

// NewTransactionRepo creates a new instance of transactionRepo.
func NewTransactionRepo(db database.IDatabase) *TransactionRepo {
	return &TransactionRepo{db: db}
}

// Create creates a new transaction record in the database.
func (r *TransactionRepo) Create(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *TransactionRepo) GetByIDAndWalletID(ctx context.Context, id, walletID uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.WithContext(ctx).Where("id = ? AND wallet_id = ?", id, walletID).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.NotFound, fmt.Errorf("transaction with ID %s not found", id))
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.NotFound, fmt.Errorf("transaction with ID %s not found", id))
		}
		return nil, err
	}
	return &transaction, err
}

func (r *TransactionRepo) Update(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *TransactionRepo) GetByWalletID(ctx context.Context, walletID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).Where("wallet_id = ?", walletID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
