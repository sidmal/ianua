package handler

import (
	"bytes"
	"context"
	"crypto/tls"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpTransport struct {
	Transport           http.RoundTripper
	logger              *zap.Logger
	responseWaitTimeout time.Duration
}

func (m *HttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), map[string]string{"name": "HttpRequest"}, time.Now())
	req = req.WithContext(ctx)

	if m.logger == nil {
		return m.makeRequest(req)
	}

	var reqBody []byte

	if req.Body != nil {
		reqBody, _ = ioutil.ReadAll(req.Body)
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	rsp, err := m.makeRequest(req)

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

func (m *HttpTransport) makeRequest(req *http.Request) (*http.Response, error) {
	var (
		rsp *http.Response
		err error
	)

	if m.Transport != nil {
		rsp, err = m.Transport.RoundTrip(req)
	} else {
		transport := http.DefaultTransport
		transport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		transport.(*http.Transport).DisableCompression = true
		transport.(*http.Transport).IdleConnTimeout = m.responseWaitTimeout
		transport.(*http.Transport).ResponseHeaderTimeout = m.responseWaitTimeout
		transport.(*http.Transport).ExpectContinueTimeout = m.responseWaitTimeout
		transport.(*http.Transport).TLSHandshakeTimeout = m.responseWaitTimeout

		rsp, err = transport.RoundTrip(req)
	}

	return rsp, err
}
