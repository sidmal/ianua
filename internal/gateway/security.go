package gateway

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"github.com/valyala/fasttemplate"
	"hash"
)

const (
	MethodNone   = "none"
	MethodMD5    = "md5"
	MethodSHA1   = "sha1"
	MethodSHA256 = "sha256"
	MethodSHA512 = "sha512"
	MethodJWT    = "jwt"
	MethodRSA    = "rsa"
)

var (
	signConstructors = map[string]func() hash.Hash{
		MethodMD5:    md5.New,
		MethodSHA1:   sha1.New,
		MethodSHA256: sha256.New,
		MethodSHA512: sha512.New,
	}
)

type Signer struct {
	Template string
	Func     func() hash.Hash
	//Methods[]
}

func NewSigner(method, template string) (*Signer, error) {
	fn, ok := signConstructors[method]
	if !ok {

	}

	signer := &Signer{
		Template: template,
		Func:     fn,
	}
	return signer, nil
}

func (m *Signer) GetSign(tmpl *fasttemplate.Template, params map[string]interface{}) []byte {
	h := m.Func()
	h.Write([]byte(tmpl.ExecuteString(params)))
	return base64.StdEncoding.EncodeToString(hash[:])
}
