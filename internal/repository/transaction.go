package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sidmal/ianua/pkg"
	"go.uber.org/zap"
)

const (
	TransactionStatusNew        = "new"
	TransactionStatusInProgress = "in_progress"
	TransactionStatusCompleted  = "completed"
	TransactionStatusRejected   = "rejected"
)

type Transaction struct {
	Model
	// The client unique identifier in billing system.
	ClientId uint64 `db:"client_id"`
	// The client name.
	ClientName string `db:"client_name"`
	// The provider unique identifier in billing system.
	// It's identifier of provider into which sending transaction.
	ProviderId uint64 `db:"provider_id"`
	// The provider name.
	ProviderName string `db:"provider_name"`
	// The service unique identifier in billing system.
	// It's identifier of service into which sending transaction.
	ServiceId uint64 `db:"service_id"`
	// The service name.
	ServiceName string `db:"service_name"`
	// The technical identifier of provider gateway handler.
	ProviderHandlerId string `db:"provider_handler_id"`
	// The transaction unique identifier in client billing system.
	ClientTxnId *string `db:"client_txn_id"`
	// The transaction unique identifier in provider billing system.
	ProviderTxnId *string `db:"provider_txn_id"`
	// The customer's account in service into which sending payment amount.
	Account string `db:"account"`
	// The key-value object that you can attach to the payment request.
	// It can be useful for storing additional information about your customer’s payment.
	Metadata map[string]interface{} `db:"metadata"`
	// The description for payment to show to customer in payment details in account statement.
	Description string `db:"description"`
	// The payment amount which was received from the client.
	IncomeAmount float32 `db:"income_amount"`
	// The currency of client's balance.
	IncomeCurrency string `db:"income_currency"`
	// The fee amount from client for payment in currency which was received from the client.
	ClientFeeInIncomeCurrency float32 `db:"client_fee_in_income_currency"`
	// The fee amount from customer for payment in currency which was received from the client.
	CustomerFeeInIncomeCurrency float32 `db:"customer_fee_in_income_currency"`
	// The payment amount which will be send to provider.
	OutcomeAmount float32 `db:"outcome_amount"`
	// The provider's currency.
	OutcomeCurrency string `db:"outcome_currency"`
	// The fee amount from client for payment in currency which will be send payment to provider.
	ClientFeeInOutcomeCurrency float32 `db:"client_fee_in_outcome_currency"`
	// The fee amount from customer for payment in currency which will be send payment to provider.
	CustomerFeeInOutcomeCurrency float32 `db:"customer_fee_in_outcome_currency"`
	// The payment amount in system accounting currency.
	AccountingAmount float32 `db:"accounting_amount"`
	// The system accounting currency.
	AccountingCurrency string `db:"accounting_currency"`
	// The fee amount from client for payment in system accounting currency.
	ClientFeeInAccountingCurrency float32 `db:"client_fee_in_accounting_currency"`
	// The fee amount from customer for payment in system accounting currency.
	CustomerFeeInAccountingCurrency float32 `db:"customer_fee_in_accounting_currency"`
	// The conversion rate value from income currency to outcome currency.
	IncomeToOutcomeRate float32 `db:"income_to_outcome_rate"`
	// The conversion rate value from income currency to accounting currency.
	IncomeToAccountingRate float32 `db:"income_to_accounting_rate"`
	// The conversion rate value from outcome currency to accounting currency.
	OutcomeToAccountingRate float32 `db:"outcome_to_accounting_rate"`
	// The transaction reject reason.
	GatewayRejectReason string `db:"gateway_reject_reason"`
	// The transaction status.
	Status string `db:"status"`
	// The client balance before transaction
	ClientBalanceBefore float32 `db:"client_balance_before"`
	// The client balance after transaction
	ClientBalanceAfter float32 `db:"client_balance_after"`
}

type transactionRepository repository

func newTransactionRepository(
	db *sqlx.DB,
	logger *zap.Logger,
) TransactionRepositoryInterface {
	repository := &transactionRepository{
		db:     db,
		logger: logger,
	}
	return repository
}

func (m *transactionRepository) GetTransactionByClientTxnId(
	ctx context.Context,
	clientId uint64,
	clientTxnId string,
) (*Transaction, error) {
	transaction := new(Transaction)
	query := "SELECT client_id, client_name, provider_id, provider_name, service_id, service_name, provider_handler_id, " +
		"client_txn_id, provider_txn_id, account, metadata, description, income_amount, income_currency, " +
		"client_fee_in_income_currency, customer_fee_in_income_currency, outcome_amount, outcome_currency, " +
		"client_fee_in_outcome_currency, customer_fee_in_outcome_currency, accounting_amount, accounting_currency, " +
		"client_fee_in_accounting_currency, customer_fee_in_accounting_currency, income_to_outcome_rate, income_to_accounting, " +
		"outcome_to_accounting, gateway_reject_reason, status, balance_before, balance_after, created_at, updated_at, deleted_at " +
		"FROM transactions WHERE client_id = $1 AND client_txn_id = $2 AND deleted_at IS NULL"
	args := []interface{}{clientId, clientTxnId}
	err := m.db.GetContext(ctx, transaction, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		m.logger.Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldFilter, query),
			zap.Any(pkg.ErrorDatabaseFieldArguments, args),
		)
		return nil, err
	}

	return transaction, nil
}

func (m *transactionRepository) Create(ctx context.Context, in *Transaction) (*Transaction, error) {
	opts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	txn, err := m.db.BeginTx(ctx, opts)

	if err != nil {
		return nil, err
	}

	clientBalanceBefore := float32(0)
	clientBalanceAfter := float32(0)
	clientFeeInIncomeCurrency := float32(0)
	ClientFeeInOutcomeCurrency := float32(0)
	ClientFeeInAccountingCurrency := float32(0)

	query := "SELECT fee_percent, balance FROM client WHERE id = $1"
	err = txn.QueryRow("SELECT $1 * (fee_percent/100), balance - ($2 + $3 * (balance/100)), balance FROM leo_paymethods WHERE id = ?",
		in.InAmount, in.InAmount, in.InAmount, in.Client.Id).
		Scan(&commission, &balanceAfter, &balance)
}

func (m *transactionRepository) Complete(txn *Transaction) error {

}

func (m *transactionRepository) Reject(txn *Transaction) error {

}
