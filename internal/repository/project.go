package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"ianua/pkg"
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
	// The global fee in percent for payment to services of service provider
	FeePercent float32 `db:"fee_percent" json:"fee_percent" validate:"omitempty,numeric,gte=0,lte=100"`
}

type ProviderServiceMap struct {
	ProviderId uint64 `db:"provider_id"`
	ServiceId  uint64 `db:"service_id"`
}

type Service struct {
	Id            primitive.ObjectID `bson:"_id"`
	Uuid          string             `bson:"uuid" json:"id"`
	Provider      *ProjectProvider   `bson:"uuid" json:"provider"`
	Name          string             `bson:"uuid" json:"name"`
	AccountRegexp string             `bson:"uuid" json:"account_regexp"`
	MinAmount     float64            `bson:"uuid" json:"min_amount"`
	MaxAmount     float64            `bson:"uuid" json:"max_amount"`
	Currency      string             `bson:"uuid" json:"currency" validate:"required,alpha,len=3"`
	FeePercent    float64            `bson:"uuid" json:"fee_percent" validate:"omitempty,numeric,gte=0,lte=100"`
	Deleted       bool               `bson:"deleted" json:"deleted"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`

	CacheAccountRegexp *regexp.Regexp `bson:"-" json:"-"`
}

type projectRepository struct {
	*repository
	cacheLifetimeProvider int
	cacheProvider         Cached
}

func newProjectRepository(
	db mongodb.SourceInterface,
	cacheLifetimeProject int,
	cacheLifetimeProvider int,
	logger *zap.Logger,
) ProjectRepositoryInterface {
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

func (m *projectRepository) GetProject(ctx context.Context, uuid string) (*Project, error) {
	cache, ok := m.cache[uuid]
	current := time.Now()

	if ok && cache.expire.After(current) {
		return cache.value.(*Project), nil
	}

	project := new(Project)
	filter := bson.M{"uuid": uuid}
	err := m.db.Collection(collectionProject).FindOne(ctx, filter).Decode(project)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, pkg.ErrorProjectNotFound
		}

		m.logger.Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionProject),
			zap.Any(pkg.ErrorDatabaseFieldFilter, filter),
		)
		return nil, pkg.ErrorUnknown
	}

	provider, err := m.GetProvider(ctx, project.Provider.ProviderId)

	if err != nil {
		return nil, err
	}

	if project.AccountRegexp != "" {
		project.CacheAccountRegexp = regexp.MustCompile(project.AccountRegexp)
	}

	if project.Currency == "" {
		project.Currency = provider.Currency
	}

	if project.FeePercent == 0 && provider.FeePercent > 0 {
		project.FeePercent = provider.FeePercent
	}

	if m.cacheLifetime > 0 {
		m.mx.Lock()
		m.cache[uuid] = &CachedValue{
			value:  project,
			expire: current.Add(time.Duration(m.cacheLifetime) * time.Second),
		}
		m.mx.Unlock()
	}

	return project, nil
}

func (m *projectRepository) GetProvider(ctx context.Context, oid primitive.ObjectID) (*Provider, error) {
	cacheKey := oid.Hex()
	cache, ok := m.cacheProvider[cacheKey]
	current := time.Now()

	if ok && cache.expire.After(current) {
		return cache.value.(*Provider), nil
	}

	provider := new(Provider)
	filter := bson.M{"_id": oid}
	err := m.db.Collection(collectionProvider).FindOne(ctx, filter).Decode(provider)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, pkg.ErrorProviderNotFound
		}

		m.logger.Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionProvider),
			zap.Any(pkg.ErrorDatabaseFieldFilter, filter),
		)
		return nil, pkg.ErrorUnknown
	}

	if m.cacheLifetime > 0 {
		m.mx.Lock()
		m.cache[cacheKey] = &CachedValue{
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
