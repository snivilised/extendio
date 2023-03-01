package log

import (
	"go.uber.org/zap"
)

// add new zap fields as and when they are required

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}
