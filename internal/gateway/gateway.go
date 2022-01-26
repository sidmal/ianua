package gateway

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"github.com/sidmal/ianua/internal/handler"
	"net/http"
	"time"
)

type Gateway struct {
	Actions []*Action
}

type Action struct {
	// Http method to request to gateway API endpoint
	Method string
	// Gateway API endpoint
	Endpoint string
	// Body template with placeholders to request to API endpoint
	Body string
	// Number of seconds to wait response from gateway API endpoint
	ResponseWaitTimeout time.Duration
	SignSettings        *SignatureSettings
}

func (m *Action) newHttpClient() *http.Client {
	httpClient := &http.Client{
		Timeout: m.ResponseWaitTimeout,
	}

	clientKey, err := base64.StdEncoding.DecodeString(cfg.Elecsnet.Login)

	if err != nil {
		return nil, err
	}

	clientCert, err := base64.StdEncoding.DecodeString(cfg.Elecsnet.Password)

	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(clientCert, clientKey)

	if err != nil {
		return nil, err
	}

	caCert, err := base64.StdEncoding.DecodeString(cfg.Elecsnet.RsaKey)

	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
		Renegotiation:      tls.RenegotiateOnceAsClient,
	}
	transport := http.DefaultTransport
	transport.(*http.Transport).TLSClientConfig = tlsConfig
	transport.(*http.Transport).TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}

	httpClient.Transport = &handler.httpTransport{
		Transport: transport,
		logger:    logger,
	}
}
