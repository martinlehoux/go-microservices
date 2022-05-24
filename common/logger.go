package common

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
	// With adds a metadata for all logs in the same context
	With(ctx context.Context, key string, value string) Logger
	// Init creates a new context with the metadata initialized
	Init(ctx context.Context) context.Context
}

// LogrusLogger is a Logger implementation that uses logrus
type LogrusLogger struct {
	logger *logrus.Logger
}

// NewLogrusLogger creates a new LogrusLogger ready for production usage
func NewLogrusLogger() LogrusLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return LogrusLogger{logger: logger}
}

func (logger *LogrusLogger) With(ctx context.Context, key string, value string) Logger {
	meta, ok := getMetaFromContext(ctx)
	if ok {
		meta[key] = value
	}
	return logger
}

func (logger *LogrusLogger) Init(ctx context.Context) context.Context {
	return context.WithValue(ctx, metaKey{}, make(meta))
}

func (logger *LogrusLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logger.getEntryFromContext(ctx).Infof(msg, args...)
}

func (logger *LogrusLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logger.getEntryFromContext(ctx).Warnf(msg, args...)
}

func (logger *LogrusLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logger.getEntryFromContext(ctx).Errorf(msg, args...)
}

type metaKey struct{}

type meta map[string]string

func getMetaFromContext(ctx context.Context) (meta, bool) {
	m, ok := ctx.Value(metaKey{}).(meta)
	return m, ok
}

func (logger *LogrusLogger) getEntryFromContext(ctx context.Context) *logrus.Entry {
	entry := logrus.NewEntry(logger.logger)
	m, ok := getMetaFromContext(ctx)
	if ok {
		for key, value := range m {
			entry = entry.WithField(key, value)
		}
	}
	return entry
}
