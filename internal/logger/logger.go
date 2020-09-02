package logger

import (
	"go.uber.org/zap"
	"sync"
)

type Logger struct {
	registry map[string]*zap.Logger
	mx       sync.Mutex
}
