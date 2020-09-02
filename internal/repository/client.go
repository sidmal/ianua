package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"ianua/pkg"
	"time"
)

type Client struct {
	Id         string   `json:"id"`
	Name       string   `json:"name" validate:"required"`
	SecretKey  string   `json:"secret_key" validate:"required,max=255"`
	FeePercent float64  `json:"fee_percent" validate:"omitempty,numeric,gte=0,lte=100"`
	Projects   []string `json:"projects" validate:"omitempty,dive,hexadecimal,len=24"`
	*ClientBalance
}

type ClientBalance struct {
	Balance  float64 `json:"balance" validate:"omitempty,numeric,gte=0"`
	Currency string  `json:"currency" validate:"required,alpha,len=3"`
}

type clientRepository courseRepository

func newClientRepository(
	db pgxpool.Pool,
	cacheLifetime int,
	logger *zap.Logger,
) ClientRepositoryInterface {
	repository := &clientRepository{
		db:            db,
		logger:        logger,
		cacheLifetime: cacheLifetime,
		cache:         make(Cached),
	}
	return repository
}

func (m *clientRepository) GetClient(ctx context.Context, uuid string) (*Client, error) {
	cache, ok := m.cache[uuid]
	current := time.Now()

	if ok && cache.expire.After(current) {
		return cache.value.(*Client), nil
	}

	client, err := m.getEntity(ctx, uuid, new(Client))

	if err != nil {
		return nil, err
	}

	if m.cacheLifetime > 0 {
		m.mx.Lock()
		m.cache[uuid] = &CachedValue{
			value:  client,
			expire: current.Add(time.Duration(m.cacheLifetime) * time.Second),
		}
		m.mx.Unlock()
	}

	return client.(*Client), nil
}

func (m *clientRepository) GetClientBalance(ctx context.Context, uuid string) (*ClientBalance, error) {
	balance, err := m.getEntity(ctx, uuid, new(ClientBalance))

	if err != nil {
		return nil, err
	}

	return balance.(*ClientBalance), nil
}

func (m *clientRepository) getEntity(ctx context.Context, uuid string) (interface{}, error) {
	client := new(Client)
	query := `select id, name, secret_key, fee_percent, projects, balance, currency 
        from ` + tableClient + ` where uuid = $1`
	err := m.db.QueryRow(ctx, query, uuid).
		Scan(client.Id, client.Name, client.SecretKey, client.FeePercent, client.Projects, client.Balance, client.Currency)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, pkg.ErrorClientNotFound
		}

		m.logger.Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionClient),
			zap.Any(pkg.ErrorDatabaseFieldFilter, filter),
		)
		return nil, pkg.ErrorUnknown
	}

	return receiver, nil
}
