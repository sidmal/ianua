package gateway

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpTransport struct {
	Transport http.RoundTripper
	logger    *zap.Logger
}

func (m *HttpSettings) getHttpTransport() (http.RoundTripper, error) {
	transport := http.DefaultTransport
	transport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport.(*http.Transport).DisableCompression = true
	transport.(*http.Transport).IdleConnTimeout = m.ResponseWaitTimeout
	transport.(*http.Transport).ResponseHeaderTimeout = m.ResponseWaitTimeout
	transport.(*http.Transport).ExpectContinueTimeout = m.ResponseWaitTimeout
	transport.(*http.Transport).TLSHandshakeTimeout = m.ResponseWaitTimeout

	if m.TLS == nil {
		return transport, nil
	}

	clientKey, err := base64.StdEncoding.DecodeString(m.TLS.ClientKey)

	if err != nil {
		return nil, err
	}

	clientCert, err := base64.StdEncoding.DecodeString(m.TLS.ClientCert)

	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(clientCert, clientKey)

	if err != nil {
		return nil, err
	}

	caCert, err := base64.StdEncoding.DecodeString(m.TLS.CaCert)

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
	transport.(*http.Transport).TLSClientConfig = tlsConfig
	transport.(*http.Transport).TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}

	return transport, nil
}

func (m *HttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), map[string]string{"name": "HttpRequest"}, time.Now())
	req = req.WithContext(ctx)

	if m.logger == nil {
		return m.Transport.RoundTrip(req)
	}

	var reqBody []byte

	if req.Body != nil {
		reqBody, _ = ioutil.ReadAll(req.Body)
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	rsp, err := m.Transport.RoundTrip(req)

	if err != nil {
		return rsp, err
	}

	var rspBody []byte

	if rsp.Body != nil {
		rspBody, err = ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
	}

	rsp.Body = ioutil.NopCloser(bytes.NewBuffer(rspBody))

	m.logger.Info(
		req.URL.String(),
		zap.String("request_method", req.Method),
		zap.Any("request_headers", req.Header),
		zap.ByteString("request_body", reqBody),
		zap.Int("response_status", rsp.StatusCode),
		zap.Any("response_headers", rsp.Header),
		zap.ByteString("response_body", rspBody),
	)

	return rsp, err
}
