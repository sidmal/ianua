package entity

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	GatewaySecurityTypeNone = "none"
	GatewaySecurityTypeHash = "hash"
	GatewaySecurityTypeJWT  = "jwt"
)

const (
	GatewaySecurityHashAlgoMD5    = "md5"
	GatewaySecurityHashAlgoSHA1   = "sha1"
	GatewaySecurityHashAlgoSHA256 = "sha256"
	GatewaySecurityHashAlgoSHA512 = "sha512"
)

const (
	HashAfterFuncAlgoHex      = "hex"
	HashAfterFuncAlgoBase64   = "base64"
	HashAfterFuncAlgoRsaPkcs1 = "rsa_pkcs1"
)

const (
	GatewaySecurityJWTOptsRspFormatRaw  = "raw"
	GatewaySecurityJWTOptsRspFormatJSON = "json"
	GatewaySecurityJWTOptsRspFormatXML  = "xml"
)

const (
	ErrorHashUnknownAlgo     = `unknown hash algo with name "%s"`
	ErrorUnknownSecurityType = `unknown security type with name "%s"`
)

var (
	HashAfterFns = map[string]func(hashed []byte, hashAlgo crypto.Hash, opts *GatewaySecurityHashAfterFuncOpts) ([]byte, error){
		HashAfterFuncAlgoHex: func(hashed []byte, _ crypto.Hash, _ *GatewaySecurityHashAfterFuncOpts) ([]byte, error) {
			return []byte(hex.EncodeToString(hashed)), nil
		},
		HashAfterFuncAlgoBase64: func(hashed []byte, _ crypto.Hash, _ *GatewaySecurityHashAfterFuncOpts) ([]byte, error) {
			return []byte(base64.StdEncoding.EncodeToString(hashed)), nil
		},
		HashAfterFuncAlgoRsaPkcs1: func(hashed []byte, hashAlgo crypto.Hash, opts *GatewaySecurityHashAfterFuncOpts) ([]byte, error) {
			return rsa.SignPKCS1v15(rand.Reader, opts.PrivateKey, hashAlgo, hashed[:])
		},
	}
)

type GatewayOpts struct {
	Name       string
	HttpClOpts *HttpClientOpts
	Security   *GatewaySecurity
	Methods    []*Method
}

type HttpClientOpts struct {
	TLS                 *TLS
	ResponseWaitTimeout time.Duration
}

type TLS struct {
	ClientKey  string
	ClientCert string
	CaCert     string
}

type GatewaySecurity struct {
	Type     string
	JWTOpts  *GatewaySecurityJWTOpts
	HashOpts *GatewaySecurityHashOpts
}

type GatewaySecurityJWTOpts struct {
	*Request
	TokenLifetime time.Duration
}

type GatewaySecurityHashOpts struct {
	Algo      string
	AfterFunc []*GatewaySecurityHashAfterFunc
}

type GatewaySecurityHashAfterFunc struct {
	Algo string
	Opts *GatewaySecurityHashAfterFuncOpts
}

type GatewaySecurityHashAfterFuncOpts struct {
	PrivateKey *rsa.PrivateKey
}

func (m *GatewaySecurityHashAfterFunc) ExecuteFunc(hashed []byte, hashAlgo crypto.Hash) ([]byte, error) {
	fn, ok := HashAfterFns[m.Algo]
	if !ok {
		return nil, fmt.Errorf(ErrorHashUnknownAlgo, m.Algo)
	}

	return fn(hashed, hashAlgo, m.Opts)
}
