package helper

import (
	"github.com/valyala/fasttemplate"
)

const (
	templateEngineStartTag = "{{"
	templateEngineEndTag   = "}}"
)

func ExecuteTemplate(template string, params map[string]interface{}) string {
	return fasttemplate.ExecuteString(template, templateEngineStartTag, templateEngineEndTag, params)
}
