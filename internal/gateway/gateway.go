package gateway

import (
	"github.com/sidmal/ianua/internal/entity"
	"github.com/sidmal/ianua/internal/gateway/signature"
	"go.uber.org/zap"
	"net/http"
)

type Gateway struct {
	HttpClient *http.Client
	Actions    []*Action
}

type Action struct {
	// Http method to request to gateway API endpoint
	Method string
	// Gateway API endpoint
	Endpoint string
	// Body template with placeholders to request to API endpoint
	Body      string
	Signature signature.Signer
}

type HttpClientSettings struct {
}

type Gateways map[string]*Gateway

func BuildGateway(gw *entity.Gateway) error {

}

func newHttpClient(opts *entity.HttpClientOpts, logger *zap.Logger) (*http.Client, error) {
	transport, err := m.getHttpTransport()
	if err != nil {
		return nil, err
	}

	cl := &http.Client{
		Timeout: m.ResponseWaitTimeout,
		Transport: &HttpTransport{
			Transport: transport,
			logger:    logger,
		},
	}

	return cl, nil
}
