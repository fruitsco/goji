package x

import "go.uber.org/zap"

func NamedLogger(name string) func(log *zap.Logger) *zap.Logger {
	return func(log *zap.Logger) *zap.Logger {
		return log.Named(name)
	}
}
