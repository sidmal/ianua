package signature

import "github.com/sidmal/ianua/internal/handler"

const (
	MethodNone = "none"
	MethodMD5  = "md5"
	MethodJWT  = "jwt"
	MethodRSA  = "rsa"
	MethodTLS  = "tls"
)

type Signer interface {
	GetMethodName() string
	ConfigureHttpTransport(tr *handler.HttpTransport) (*handler.HttpTransport, error)
}
