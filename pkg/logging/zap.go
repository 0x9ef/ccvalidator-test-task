package logging

import (
	"context"
	"strings"

	zapLib "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zap struct {
	logger *zapLib.SugaredLogger
}

// Check if implements the interface.
var _ Logger = (*zap)(nil)

func NewZap(level string) *zap {
	var l zapcore.Level
	switch strings.ToLower(level) {
	case "error":
		l = zapcore.ErrorLevel
	case "warn":
		l = zapcore.WarnLevel
	case "info":
		l = zapcore.InfoLevel
	case "debug":
		l = zapcore.DebugLevel
	default:
		l = zapcore.InfoLevel
	}

	config := zapLib.Config{
		Development:      false,
		Encoding:         "json",
		Level:            zapLib.NewAtomicLevelAt(l),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeDuration: zapcore.SecondsDurationEncoder,
			LevelKey:       "severity",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			TimeKey:        "timestamp",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			NameKey:        "name",
			EncodeName:     zapcore.FullNameEncoder,
			MessageKey:     "message",
			StacktraceKey:  "",
			LineEnding:     "\n",
		},
	}

	logger, _ := config.Build()
	logger = logger.WithOptions(zapLib.AddCallerSkip(1)) // trace last call frame

	return &zap{logger: logger.Sugar()}
}

func (l *zap) Named(name string) Logger {
	return &zap{
		logger: l.logger.Named(name),
	}
}

func (l *zap) With(args ...interface{}) Logger {
	return &zap{
		logger: l.logger.With(args...),
	}
}

func (l *zap) WithContext(ctx context.Context, key string) Logger {
	if key == "" { // if no context key, return same logger
		return l
	}
	return &zap{
		logger: l.logger.With(key, ctx.Value(key)),
	}
}

func (l *zap) Debug(msg string, args ...interface{}) { l.logger.Debugw(msg, args...) }
func (l *zap) Info(msg string, args ...interface{})  { l.logger.Infow(msg, args...) }
func (l *zap) Error(msg string, args ...interface{}) { l.logger.Errorw(msg, args...) }
func (l *zap) Fatal(msg string, args ...interface{}) { l.logger.Fatalw(msg, args...) }
