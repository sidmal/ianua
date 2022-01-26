package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	collectionCourse      = "course"
	collectionProject     = "project"
	collectionProvider    = "provider"
	tableClient           = "client"
	collectionTransaction = "transaction"
)

type Interface interface {
	GetClientRepository() MerchantRepositoryInterface
	GetCourseRepository() CourseRepositoryInterface
	GetProjectRepository() ProviderRepositoryInterface
	GetTransactionRepository() TransactionRepositoryInterface
	GetAccountingEntryRepository() AccountingEntryRepositoryInterface
}

type CacheLifetime struct {
	Course  int
	Client  int
	Project int
}

type Repository struct {
	course  CourseRepositoryInterface
	client  MerchantRepositoryInterface
	project ProviderRepositoryInterface
	order   TransactionRepositoryInterface
}

type Cached map[string]*CachedValue

type CachedValue struct {
	value  interface{}
	expire time.Time
}

type repository struct {
	db            *sqlx.DB
	logger        *zap.Logger
	mx            sync.Mutex
	cacheLifetime int
	cache         Cached
}

type Model struct {
	Id        uint64     `db:"id"`
	Uuid      string     `db:"uuid"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Cache interface {
	GetAllCached() Cached
	RemoveCachedByKey(key interface{})
	RemoveAllCached()
}

type CourseRepositoryInterface interface {
	GetCourseRate(ctx context.Context, from, to string) (float32, error)
}

type MerchantRepositoryInterface interface {
	GetClient(ctx context.Context, uuid string) (*Client, error)
}

type ProviderRepositoryInterface interface {
	GetService(ctx context.Context, uuid string) (*Service, error)
	GetProvider(ctx context.Context, oid primitive.ObjectID) (*Provider, error)
}

type TransactionRepositoryInterface interface {
	GetTransactionByClientTxnId(ctx context.Context, clientId uint64, clientTxnId string) (*Transaction, error)
	//CreateOrder(in *PaymentIn) (*Order, error)
	//Update(id int64, status int, gatewayOrderId string) error
	//RejectPayment(id int64) error
}

func NewRepository(db *sqlx.DB, cacheLifetime *CacheLifetime, logger *zap.Logger) Interface {
	repository := &Repository{
		course: newCourseRepository(db, cacheLifetime.Course, logger),
	}

	return repository
}

func (m *Repository) GetClientRepository() MerchantRepositoryInterface {
	return m.client
}

func (m *Repository) GetCourseRepository() CourseRepositoryInterface {
	return m.course
}

func (m *Repository) GetProjectRepository() ProviderRepositoryInterface {
	return m.project
}

func (m *Repository) GetOrderRepository() OrderRepositoryInterface {
	return m.order
}

func (m *Repository) GetAccountingEntryRepository() AccountingEntryRepositoryInterface {
	return m.accountingEntry
}
