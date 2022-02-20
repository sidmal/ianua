package entity

import "time"

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
	GatewaySecurityHashAfterFuncAlgoRsaPkcs1 = "rsa_pkcs1"
	GatewaySecurityHashAfterFuncAlgoBase64   = "base64"
)

type Gateway struct {
	Name       string
	HttpClOpts *HttpClientOpts
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
	Type string
}

type GatewaySecurityJWTOpts struct {
	Url                    string
	Headers                []map[string]string
	RequestMethod          string
	RequestBody            string
	ResponseTokenFieldName string
	TokenLifetime          time.Duration
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
	PrivateKey string
}
