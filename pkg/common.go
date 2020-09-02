package pkg

const (
	RateMultiplier   = 1000000
	AmountMultiplier = 100
)

type BaseRequest struct {
	Account   string `json:"account" validate:"required,min=1"`
	ProjectId string `json:"project_id" validate:"required,uuid"`
}

type StatusRequest struct {
	OrderId string `json:"order_id" validate:"required"`
}

type PaymentRequest struct {
	BaseRequest
	StatusRequest
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Amount   float32                `json:"amount" validate:"required,gt=0"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func NewError(code, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

func (m *Error) SetDetails(details string) *Error {
	err := &Error{
		Code:    m.Code,
		Message: m.Message,
		Details: details,
	}
	return err
}

func (m *Error) Error() string {
	return m.Message
}

const (
	ErrorDatabaseQueryFailed    = "query to database collection failed"
	ErrorDatabaseFieldFilter    = "query"
	ErrorDatabaseFieldArguments = "query"
)

var (
	ErrorMerchantNotFound = NewError("mr000008", "client with specified identifier not found")
	ErrorServiceNotFound  = NewError("mr000008", "service with specified identifier not found")
	ErrorServiceInactive  = NewError("mr000008", "service with specified identifier is inactive")
	ErrorProviderNotFound = NewError("mr000008", "provider for project with received identifier not found")
	ErrorProviderInactive = NewError("mr000008", "provider for project with received identifier is inactive")
	ErrorCourseNotFound   = NewError("mr000008", "rate for currency conversion from client balance currency to project recipient currency not found")
	ErrorUnknown          = NewError("mr000008", "unknown error, try request later")
)
