package transport

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/sidmal/ianua/internal/entity"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	errorRspNotSuccessCode = `request failed with reason in body: "%s"`
	errorRspFormatUnknown  = `response format "%s" is unknown`
	errorRspFieldNotObject = `field with name "%s" not contain children nodes`
	errorRspFieldNotFound  = `field with name "%s" not found in response`
)

type HttpClient struct {
	httpCl *http.Client
}

type HttpTransport struct {
	Transport http.RoundTripper
	logger    *zap.Logger
}

func NewHttpClient(opts *entity.HttpClientOpts, logger *zap.Logger) (*HttpClient, error) {
	transport, err := getHttpTransport(opts)
	if err != nil {
		return nil, err
	}

	cl := &HttpClient{
		httpCl: &http.Client{
			Transport: &HttpTransport{
				Transport: transport,
				logger:    logger,
			},
			Timeout: opts.ResponseWaitTimeout,
		},
	}

	return cl, nil
}

func (m *HttpClient) MakeRequest(opts *entity.Request, reqBody string) (string, error) {
	req, err := http.NewRequest(opts.ReqMethod, opts.Url, strings.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	for key, val := range opts.Headers {
		req.Header.Add(key, val)
	}

	rsp, err := m.httpCl.Do(req)
	if err != nil {
		return "", err
	}

	rspBody, err := ioutil.ReadAll(rsp.Body)
	_ = rsp.Body.Close()
	if err != nil {
		return "", err
	}

	if rsp.StatusCode != opts.RspSuccessCode {
		return "", fmt.Errorf(errorRspNotSuccessCode, string(rspBody))
	}

	if len(opts.RspResultFieldNames) == 0 || opts.RspFormat == entity.GatewaySecurityJWTOptsRspFormatRaw {
		return string(rspBody), nil
	}

	return m.getResponseResultFieldValue(rspBody, opts.RspFormat, opts.RspResultFieldNames)
}

func (m *HttpClient) getResponseResultFieldValue(
	rspBody []byte,
	rspFormat string,
	rspResultFieldNames []string,
) (string, error) {
	var (
		res = make(map[string]interface{})
		err error
	)

	switch rspFormat {
	case entity.GatewaySecurityJWTOptsRspFormatJSON:
		err = json.Unmarshal(rspBody, &res)
	case entity.GatewaySecurityJWTOptsRspFormatXML:
		err = xml.Unmarshal(rspBody, &res)
	default:
		return "", fmt.Errorf(errorRspFormatUnknown, rspFormat)
	}

	if err != nil {
		return "", err
	}

	resFieldVal := res[rspResultFieldNames[0]]
	rspResultFieldNamesLen := len(rspResultFieldNames)

	if rspResultFieldNamesLen == 1 {
		return fmt.Sprint(resFieldVal), nil
	}

	for i := 1; i < rspResultFieldNamesLen; i++ {
		resFieldValTyped, ok := resFieldVal.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf(errorRspFieldNotObject, rspResultFieldNames[i-1])
		}

		resFieldVal, ok = resFieldValTyped[rspResultFieldNames[i]]
		if !ok {
			return "", fmt.Errorf(errorRspFieldNotFound, rspResultFieldNames[i])
		}
	}

	return fmt.Sprint(resFieldVal), nil
}

func getHttpTransport(opts *entity.HttpClientOpts) (http.RoundTripper, error) {
	transport := http.DefaultTransport
	transport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport.(*http.Transport).DisableCompression = true
	transport.(*http.Transport).IdleConnTimeout = opts.ResponseWaitTimeout
	transport.(*http.Transport).ResponseHeaderTimeout = opts.ResponseWaitTimeout
	transport.(*http.Transport).ExpectContinueTimeout = opts.ResponseWaitTimeout
	transport.(*http.Transport).TLSHandshakeTimeout = opts.ResponseWaitTimeout

	if opts.TLS == nil {
		return transport, nil
	}

	clientKey, err := base64.StdEncoding.DecodeString(opts.TLS.ClientKey)

	if err != nil {
		return nil, err
	}

	clientCert, err := base64.StdEncoding.DecodeString(opts.TLS.ClientCert)

	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(clientCert, clientKey)

	if err != nil {
		return nil, err
	}

	caCert, err := base64.StdEncoding.DecodeString(opts.TLS.CaCert)

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
