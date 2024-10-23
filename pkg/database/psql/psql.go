package psql

import (
	"context"

	db "github.com/3bd-dev/wallet-service/pkg/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

type Config struct {
	URL          string
	MaxOpenConns int
	MaxIdleConns int
}

func Open(cfg Config) (db.IDatabase, error) {
	orm, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.URL,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sql, err := orm.DB()
	if err != nil {
		return nil, err
	}
	sql.SetMaxOpenConns(cfg.MaxOpenConns)
	sql.SetMaxIdleConns(cfg.MaxIdleConns)

	return &database{db: orm}, nil
}

func (d *database) Ping() error {
	sql, err := d.db.DB()
	if err != nil {
		return err
	}

	return sql.Ping()
}

func (d *database) Close() error {
	sql, err := d.db.DB()
	if err != nil {
		return err
	}

	return sql.Close()
}

func (d *database) Client() db.IDatabase {
	return d
}

func (d *database) WithContext(ctx context.Context) *gorm.DB {
	tx := ctx.Value(db.ContextKeyDBTx)
	if tx == nil {
		return d.db.WithContext(ctx)
	}

	dbTx, ok := tx.(*database)
	if !ok {
		return d.db.WithContext(ctx)
	}

	return dbTx.db
}

func (d *database) Begin() db.IDatabase {
	return &database{
		db: d.db.Begin(),
	}
}

func (d *database) Commit() error {
	return d.db.Commit().Error
}

func (d *database) Rollback() error {
	return d.db.Rollback().Error
}
