package gateway

import "github.com/sidmal/ianua/internal/entity"

type Method struct {
	*entity.Method
}

func NewMethod(method *entity.Method) *Method {
	return &Method{
		Method: method,
	}
}

func (m *Method) ExecuteMethod(params map[string]interface{}) error {

}
