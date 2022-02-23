package security

import (
	"github.com/sidmal/ianua/internal/entity"
	"time"
)

type JWT struct {
	sign  *Sign
	opts  *entity.GatewaySecurityJWTOpts
	token *token
}

type token struct {
	token  string
	expire time.Time
}

func newJWTSigner(sign *Sign, opts *entity.GatewaySecurityJWTOpts) Signer {
	return &JWT{
		sign: sign,
		opts: opts,
	}
}

func (m *JWT) GetSignature(_ string, _ map[string]interface{}) (string, error) {
	if m.token != nil && m.token.expire.After(time.Now()) {
		return m.token.token, nil
	}

	res, err := m.sign.httpCl.MakeRequest(m.opts.Request, m.opts.ReqBody)
	if err != nil {
		return "", err
	}

	m.token = &token{
		token:  res,
		expire: time.Now().Add(m.opts.TokenLifetime),
	}

	return m.token.token, nil
}
