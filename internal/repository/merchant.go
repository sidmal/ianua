package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sidmal/ianua/pkg"
	"go.uber.org/zap"
	"time"
)

type Merchant struct {
	Model
	Name       string   `db:"name" json:"name" validate:"required"`
	SecretKey  string   `db:"secret_key" json:"secret_key" validate:"required,max=255"`
	FeePercent float64  `db:"fee_percent" json:"fee_percent" validate:"omitempty,numeric,gte=0,lte=100"`
	Balance    float64  `db:"balance" json:"balance" validate:"omitempty,numeric,gte=0"`
	Currency   string   `db:"currency" json:"currency" validate:"required,alpha,len=3"`
	Projects   []string `db:"-" json:"-"`
}

type merchantRepository courseRepository

func newMerchantRepository(
	db *sqlx.DB,
	cacheLifetime int,
	logger *zap.Logger,
) MerchantRepositoryInterface {
	repository := &merchantRepository{
		db:            db,
		logger:        logger,
		cacheLifetime: cacheLifetime,
		cache:         make(Cached),
	}
	return repository
}

func (m *merchantRepository) GetMerchant(ctx context.Context, uuid string) (*Merchant, error) {
	cache, ok := m.cache[uuid]
	current := time.Now()

	if ok && cache.expire.After(current) {
		return cache.value.(*Merchant), nil
	}

	merchant := new(Merchant)
	query := `SELECT uuid, name, secret_key, fee_percent, balance, currency FROM merchants WHERE uuid = $1 AND deleted_at IS NULL`
	args := []interface{}{uuid}
	err := m.db.GetContext(ctx, merchant, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrorMerchantNotFound
		}

		m.logger.Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldFilter, query),
			zap.Any(pkg.ErrorDatabaseFieldArguments, args),
		)
		return nil, err
	}

	if m.cacheLifetime > 0 {
		m.mx.Lock()
		m.cache[uuid] = &CachedValue{
			value:  merchant,
			expire: current.Add(time.Duration(m.cacheLifetime) * time.Second),
		}
		m.mx.Unlock()
	}

	return merchant, nil
}
