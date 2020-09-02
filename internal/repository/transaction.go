package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"ianua/pkg"
)

const (
	TransactionStatusNew        = "new"
	TransactionStatusInProgress = "in_progress"
	TransactionStatusCompleted  = "completed"
	TransactionStatusRejected   = "rejected"
)

type TransactionObject struct {
	Id   primitive.ObjectID `bson:"id"`
	Name string             `bson:"string"`
}

type TransactionProvider struct {
	TransactionObject
	ServiceId string `bson:"service_id"`
	HandlerId string `bson:"handler_id"`
}

type TransactionAmount struct {
	Amount    int64  `bson:"amount"`
	Currency  string `bson:"currency"`
	FeeClient int64  `bson:"fee_client"`
	FeeUser   int64  `bson:"fee_user"`
	Rate      int64  `bson:"rate"`
}

type TransactionFee struct {
	InIncomeCurrency     int64 `bson:"in_income_currency"`
	InOutcomeCurrency    int64 `bson:"in_outcome_currency"`
	InAccountingCurrency int64 `bson:"in_accounting_currency"`
}

type Transaction struct {
	Id                  primitive.ObjectID     `bson:"_id"`
	Uuid                string                 `bson:"uuid"`
	Client              *TransactionObject     `bson:"client"`
	Provider            *TransactionProvider   `bson:"provider"`
	Project             *TransactionObject     `bson:"project"`
	ClientTxnId         string                 `bson:"client_txn_id"`
	ProviderTxnId       string                 `bson:"provider_txn_id"`
	Account             string                 `bson:"account"`
	Metadata            map[string]interface{} `bson:"metadata"`
	Description         string                 `bson:"description"`
	Income              *TransactionAmount     `bson:"income"`
	Outcome             *TransactionAmount     `bson:"outcome"`
	Accounting          *TransactionAmount     `bson:"accounting"`
	GatewayRejectReason string                 `bson:"gateway_reject_reason"`
	Status              string                 `bson:"status"`
}

type transactionRepository repository

func newTransactionRepository(
	db mongodb.SourceInterface,
	logger *zap.Logger,
) TransactionRepositoryInterface {
	repository := &transactionRepository{
		db:     db,
		logger: logger,
	}
	return repository
}

func (m *transactionRepository) GetTransactionByClientTransactionId(
	ctx context.Context,
	clientId int,
	transactionId string,
) (*Transaction, error) {
	transaction := new(Transaction)
	filter := bson.M{
		"client.id":     clientId,
		"client_txn_id": transactionId,
	}
	err := m.db.Collection(collectionTransaction).FindOne(ctx, filter).Decode(transaction)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		m.logger.Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionTransaction),
			zap.Any(pkg.ErrorDatabaseFieldFilter, filter),
		)
		return nil, pkg.ErrorUnknown
	}

	return transaction, nil
}

func (m *transactionRepository) Create(in *pkg.PaymentRequest) (*Transaction, error) {

}

func (m *transactionRepository) Complete(txn *Transaction) error {

}

func (m *transactionRepository) Reject(txn *Transaction) error {

}
