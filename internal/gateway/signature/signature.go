package signature

import "github.com/valyala/fasttemplate"

const (
	MethodNone = "none"
	MethodMD5  = "md5"
	MethodJWT  = "jwt"
	MethodRSA  = "rsa"
)

type Signer interface {
	GetMethodNameGetSignature(tmpl *fasttemplate.Template, params map[string]interface{}) string
}
