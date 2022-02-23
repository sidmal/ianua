package gateway

import (
	"github.com/sidmal/ianua/internal/entity"
	"github.com/sidmal/ianua/internal/gateway/security"
	"github.com/sidmal/ianua/internal/gateway/transport"
	"go.uber.org/zap"
	"net/http"
)

type Gateways map[string]*Gateway

type Gateway struct {
	transport *transport.HttpClient
	Methods   []*Action
}

type Action struct {
	// Http method to request to gateway API endpoint
	Method string
	// Gateway API endpoint
	Endpoint string
	// Body template with placeholders to request to API endpoint
	Body      string
	Signature security.Signer
}

type HttpTransport struct {
	Transport http.RoundTripper
	logger    *zap.Logger
}

func NewGateway(opts *entity.GatewayOpts, logger *zap.Logger) (*Gateway, error) {
	httpCl, err := transport.NewHttpClient(opts.HttpClOpts, logger)
	if err != nil {
		return nil, err
	}

	gw := &Gateway{
		transport: httpCl,
	}

	if opts.Security != nil {

	}

}
