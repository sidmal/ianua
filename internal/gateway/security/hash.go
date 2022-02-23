package security

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/sidmal/ianua/internal/entity"
	"hash"
)

const (
	errorHashUnknownAlgo = `unknown algo method name "%s"`
)

var (
	hashFns = map[string]func() hash.Hash{
		entity.GatewaySecurityHashAlgoMD5:    md5.New,
		entity.GatewaySecurityHashAlgoSHA1:   sha1.New,
		entity.GatewaySecurityHashAlgoSHA256: sha256.New,
		entity.GatewaySecurityHashAlgoSHA512: sha512.New,
	}
)

type Hash struct {
	sign *Sign
	opts *entity.GatewaySecurityHashOpts
}

func newHashSigner(sign *Sign, opts *entity.GatewaySecurityHashOpts) Signer {
	return &Hash{
		sign: sign,
		opts: opts,
	}
}

func (m *Hash) GetSignature(template string, params map[string]interface{}) (string, error) {
	hashFn, ok := hashFns[m.opts.Algo]
	if !ok {
		return "", fmt.Errorf(errorHashUnknownAlgo, m.opts.Algo)
	}

	h := hashFn()
	h.Write([]byte(m.sign.executeTemplate(template, params)))
	hashed := h.Sum(nil)
	if len(m.opts.AfterFunc) == 0 {
		return hex.EncodeToString(hashed), nil
	}

	for _, val := range m.opts.AfterFunc {
		hashed, err := fn()
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed[:])

	return base64.StdEncoding.EncodeToString(hash[:])
}

func (m *Hash) base64() s {

}
