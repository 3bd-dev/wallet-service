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

// WalletRepo defines the interface for wallet repository.
type WalletRepo struct {
	db database.IDatabase
}

// NewWalletRepo creates a new instance of walletRepo.
func NewWalletRepo(db database.IDatabase) *WalletRepo {
	return &WalletRepo{db: db}
}

// GetByID retrieves a wallet record by its ID.
func (r *WalletRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.NotFound, fmt.Errorf("wallet with ID %s not found", id))
		}
		return nil, err
	}

	return &wallet, err
}

// Create creates a new wallet record in the database.
func (r *WalletRepo) Create(ctx context.Context, wallet *models.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

// List retrieves all wallet records from the database.
func (r *WalletRepo) List(ctx context.Context) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.WithContext(ctx).Find(&wallets).Error
	if err != nil {
		return nil, err
	}

	return wallets, nil
}
