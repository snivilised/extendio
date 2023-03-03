package log

import (
	"github.com/snivilised/extendio/xfs/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field
type Level = zapcore.Level

const (
	InfoLevel  = zapcore.InfoLevel
	ErrorLevel = zapcore.ErrorLevel
)

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Sync() error
}

type Ref utils.RoProp[Logger]

type Rotation struct {
	Filename       string
	MaxSizeInMb    int
	MaxNoOfBackups int
	MaxAgeInDays   int
}

type LoggerInfo struct {
	Rotation

	Enabled         bool
	Path            string
	TimeStampFormat string
	Level           Level
}
