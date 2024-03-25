package logging

import "context"

type Logger interface {
	Named(name string) Logger
	With(args ...interface{}) Logger
	WithContext(ctx context.Context, key string) Logger
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}
