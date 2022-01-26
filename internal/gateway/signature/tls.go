package signature

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"github.com/sidmal/ianua/internal/handler"
	"net/http"
)

type TLS struct {
	ClientKey  string
	ClientCert string
	CaCert     string
}

func (m *TLS) GetMethodName() string {
	return MethodTLS
}

func (m *TLS) ConfigureHttpTransport(tr *handler.HttpTransport) (*handler.HttpTransport, error) {
	clientKey, err := base64.StdEncoding.DecodeString(m.ClientKey)

	if err != nil {
		return nil, err
	}

	clientCert, err := base64.StdEncoding.DecodeString(m.ClientCert)

	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(clientCert, clientKey)

	if err != nil {
		return nil, err
	}

	caCert, err := base64.StdEncoding.DecodeString(m.CaCert)

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
	tr.Transport = http.DefaultTransport
	tr.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	tr.Transport.(*http.Transport).TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}

	return tr, nil
}
