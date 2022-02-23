package security

import (
	"fmt"
	"github.com/sidmal/ianua/internal/entity"
	"github.com/sidmal/ianua/internal/gateway/transport"
)

type Signer interface {
	GetSignature(template string, params map[string]interface{}) (string, error)
}

func NewSigner(httpCl *transport.HttpClient, opts *entity.GatewaySecurity) (Signer, error) {
	if opts == nil {
		return newNoneSigner(), nil
	}

	switch opts.Type {
	case entity.GatewaySecurityTypeHash:
		return newHashSigner(opts.HashOpts), nil
	case entity.GatewaySecurityTypeJWT:
		return newJWTSigner(httpCl, opts.JWTOpts), nil
	case entity.GatewaySecurityTypeNone:
		return newNoneSigner(), nil
	}

	return nil, fmt.Errorf(entity.ErrorUnknownSecurityType, opts.Type)
}
