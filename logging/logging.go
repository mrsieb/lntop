package logging

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field
type ObjectEncoder = zapcore.ObjectEncoder

type Logger interface {
	Info(string, ...Field)
	Error(string, ...Field)
	Sync() error
	Debug(string, ...Field)
	With(fields ...Field) *zap.Logger
}

func String(k, v string) Field {
	return zap.String(k, v)
}

func Duration(k string, d time.Duration) Field {
	return zap.Duration(k, d)
}

func Int(k string, i int) Field {
	return zap.Int(k, i)
}

func Int64(k string, i int64) Field {
	return zap.Int64(k, i)
}

func Error(v error) Field {
	return zap.Error(v)
}

func Object(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}

type Config struct {
	Environment string `json:"environment"`
}

func NewCliLogger(c *Config) (Logger, error) {
	cfg := zap.NewProductionConfig()
	if c.Environment == "debug" {
		cfg = zap.NewDevelopmentConfig()
	}
	return cfg.Build()
}
