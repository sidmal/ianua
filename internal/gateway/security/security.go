package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"github.com/sidmal/ianua/internal/gateway/transport"
	"github.com/valyala/fasttemplate"
)

const (
	templateEngineStartTag = "{{"
	templateEngineEndTag   = "}}"
)

type Signer interface {
	GetSignature(template string, params map[string]interface{}) (string, error)
}

type Sign struct {
	httpCl *transport.HttpClient
}

func (m *Sign) executeTemplate(template string, params map[string]interface{}) string {
	return fasttemplate.ExecuteString(template, templateEngineStartTag, templateEngineEndTag, params)
}

func (m *Sign) rsaPkcs1Base64Encrypt(privateKey *rsa.PrivateKey, hashed []byte) (string, error) {
	res, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(res), nil
}

func (m *Sign) rsaPkcs1HexEncrypt(privateKey *rsa.PrivateKey, hashed []byte) (string, error) {
	res, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(res), nil
}
