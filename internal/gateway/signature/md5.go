package signature

import (
	"github.com/sidmal/ianua/internal/handler"
)

type MD5 struct {
	Method string
}

func (m *MD5) GetMethodName() string {
	return MethodMD5
}

func (m *MD5) ConfigureHttpTransport(tr *handler.HttpTransport) (*handler.HttpTransport, error) {
	return tr, nil
}
