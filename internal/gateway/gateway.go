package gateway

import (
	"fmt"
	"github.com/sidmal/ianua/internal/entity"
	"github.com/sidmal/ianua/internal/gateway/security"
	"github.com/sidmal/ianua/internal/gateway/transport"
	"go.uber.org/zap"
	"net/http"
)

type Gateways map[string]*Gateway

type Gateway struct {
	transport *transport.HttpClient
	signer    security.Signer
	methods   []*entity.Method
}

type HttpTransport struct {
	Transport http.RoundTripper
	logger    *zap.Logger
}

var (
	gateways Gateways
)

const (
	errorGatewayNotFound = `gateway with name "%s" not found`
)

func NewGateway(opts *entity.GatewayOpts, logger *zap.Logger) (*Gateway, error) {
	httpCl, err := transport.NewHttpClient(opts.HttpClOpts, logger)
	if err != nil {
		return nil, err
	}

	signer, err := security.NewSigner(httpCl, opts.Security)
	if err != nil {
		return nil, err
	}

	gw := &Gateway{
		transport: httpCl,
		signer:    signer,
		methods:   opts.Methods,
	}
	gateways[opts.Name] = gw

	return gw, nil
}

func ExecuteGatewayMethods(gatewayName string, params map[string]interface{}) error {
	gw, ok := gateways[gatewayName]
	if !ok {
		return fmt.Errorf(`gateway with name "%s" not found`, gatewayName)
	}

	for _, method := range gw.methods {
		gw.transport.MakeRequest(method.Request)
	}
}
