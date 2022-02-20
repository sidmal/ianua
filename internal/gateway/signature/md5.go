package signature

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/valyala/fasttemplate"
)

type MD5 struct {
	// Template with placeholders to create signature
	Template string
}

func (m *MD5) GetMethodName() string {
	return MethodMD5
}

func (m *MD5) GetSignature(tmpl *fasttemplate.Template, params map[string]interface{}) string {
	str := []byte(tmpl.ExecuteString(params))
	hash := md5.Sum(str)
	return base64.StdEncoding.EncodeToString(hash[:])
}
