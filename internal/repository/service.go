package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sidmal/ianua/pkg"
	"go.uber.org/zap"
	"regexp"
	"time"
)

type Provider struct {
	Model
	// The service provider name
	Name string `db:"name" json:"name" validate:"required,min=1"`
	// The default currency to payment to service provider services
	Currency string `db:"currency" json:"currency" validate:"required,alpha,len=3"`
	// The handler technical identifier to process payments to services of service provider
	Handler string `db:"handler" json:"handler" validate:"required,min=1"`
}

type Service struct {
	Model
	// The provider identifier
	ProviderUuid string `db:"provider_uuid" json:"provider_id"`
	// The service name
	Name string `db:"name" json:"name"`
	// The account regexp for check customer's account before send request to API gateway
	AccountRegexp string `db:"account_regexp" json:"account_regexp"`
	// The phrase to get account from customer if payment init from payment form
	AccountPhrase string `db:"account_phrase" json:"account_phrase"`
	// The minimal amount in service provider currency to pay into service
	MinAmount float64 `db:"min_amount" json:"min_amount"`
	// The maximal amount in service provider currency to pay into service
	MaxAmount float64 `db:"max_amount" json:"max_amount"`
	// The unique service identifier in provider's billing system
	ExternalId string `db:"external_id" json:"external_id"`
	// Fee cost by which the payment amount must be reduced, i.e. customer receiving amount which will be reduced by this fee.
	// For example, if customer's want to pay amount 100 USD and this fee value is 3%, than after payment customer's
	// receive amount which will calculate by next formula: 100 - 100 * 0.03 = 97 USD.
	FeePercent float64 `db:"fee_percent" json:"fee_percent" validate:"omitempty,numeric,gte=0,lte=100"`
	// The technical field to save compiled regexp in cache
	CacheAccountRegexp *regexp.Regexp `db:"-" json:"-"`
	// The technical field to get provider data from service data
	Provider *Provider `db:"-" json:"-"`
}

type projectRepository struct {
	*repository
	cacheLifetimeProvider int
	cacheProvider         Cached
}

func newProviderRepository(
	db *sqlx.DB,
	cacheLifetimeProject int,
	cacheLifetimeProvider int,
	logger *zap.Logger,
) ProviderRepositoryInterface {
	repository := &projectRepository{
		repository: &repository{
			db:            db,
			logger:        logger,
			cacheLifetime: cacheLifetimeProject,
			cache:         make(Cached),
		},
		cacheLifetimeProvider: cacheLifetimeProvider,
		cacheProvider:         make(Cached),
	}
	return repository
}

func (m *projectRepository) GetService(ctx context.Context, uuid string) (*Service, error) {
	cache, ok := m.cache[uuid]
	current := time.Now()

	if ok && cache.expire.After(current) {
		return cache.value.(*Service), nil
	}

	service := new(Service)
	query := `SELECT id, uuid, provider_id, name, account_regexp, account_phrase, external_id, min_amount, max_amount, 
		fee_percent, deleted_at FROM services WHERE uuid = $1`
	args := []interface{}{uuid}
	err := m.db.GetContext(ctx, service, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrorServiceNotFound
		}

		m.logger.Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldFilter, query),
			zap.Any(pkg.ErrorDatabaseFieldArguments, args),
		)
		return nil, err
	}

	if service.DeletedAt.IsZero() {
		return nil, pkg.ErrorServiceInactive

	}

	provider, err := m.GetProvider(ctx, service.ProviderUuid)

	if err != nil {
		return nil, err
	}

	if service.AccountRegexp != "" {
		service.CacheAccountRegexp = regexp.MustCompile(service.AccountRegexp)
	}

	service.Provider = provider

	if m.cacheLifetime > 0 {
		m.mx.Lock()
		m.cache[uuid] = &CachedValue{
			value:  service,
			expire: current.Add(time.Duration(m.cacheLifetime) * time.Second),
		}
		m.mx.Unlock()
	}

	return service, nil
}

func (m *projectRepository) GetProvider(ctx context.Context, uuid string) (*Provider, error) {
	cache, ok := m.cacheProvider[uuid]
	current := time.Now()

	if ok && cache.expire.After(current) {
		return cache.value.(*Provider), nil
	}

	provider := new(Provider)
	query := `SELECT id, uuid, name, currency, handler, deleted_at FROM providers WHERE uuid = $1`
	args := []interface{}{uuid}
	err := m.db.GetContext(ctx, provider, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrorProviderNotFound
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
			value:  provider,
			expire: current.Add(time.Duration(m.cacheLifetimeProvider) * time.Second),
		}
		m.mx.Unlock()
	}

	return provider, nil
}

func (m *projectRepository) GetAllCached() Cached {
	return m.cache
}

func (m *projectRepository) RemoveCachedByKey(key interface{}) {
	m.mx.Lock()
	delete(m.cache, key)
	m.mx.Unlock()
}

func (m *projectRepository) RemoveAllCached() {
	m.mx.Lock()
	m.cache = make(Cached)
	m.mx.Unlock()
}
