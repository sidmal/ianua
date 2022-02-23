package security

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"github.com/sidmal/ianua/internal/entity"
	"github.com/sidmal/ianua/internal/gateway/helper"
)

var (
	hashFns = map[string]crypto.Hash{
		entity.GatewaySecurityHashAlgoMD5:    crypto.MD5,
		entity.GatewaySecurityHashAlgoSHA1:   crypto.SHA1,
		entity.GatewaySecurityHashAlgoSHA256: crypto.SHA256,
		entity.GatewaySecurityHashAlgoSHA512: crypto.SHA512,
	}
)

type Hash struct {
	opts *entity.GatewaySecurityHashOpts
}

func newHashSigner(opts *entity.GatewaySecurityHashOpts) Signer {
	return &Hash{
		opts: opts,
	}
}

func (m *Hash) GetSignature(template string, params map[string]interface{}) (string, error) {
	hashType, ok := hashFns[m.opts.Algo]
	if !ok {
		return "", fmt.Errorf(entity.ErrorHashUnknownAlgo, m.opts.Algo)
	}

	var (
		hashed []byte
		err    error
	)

	h := hashType.New()
	h.Write([]byte(helper.ExecuteTemplate(template, params)))
	hashed = h.Sum(nil)
	if len(m.opts.AfterFunc) == 0 {
		return hex.EncodeToString(hashed), nil
	}

	for _, val := range m.opts.AfterFunc {
		hashed, err = val.ExecuteFunc(hashed, hashType)
		if err != nil {
			return "", err
		}
	}

	return string(hashed), nil
}
