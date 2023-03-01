package log

import (
	"github.com/snivilised/extendio/xfs/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

type Flusher = zapcore.WriteSyncer

type Level = zapcore.Level

type Handle = *zap.Logger

type Ref utils.RoProp[Handle]

const (
	InfoLevel  = zapcore.InfoLevel
	DebugLevel = zapcore.DebugLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

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
