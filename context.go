package goji

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

// MARK: - Context

type contextKey int

var configKey = contextKey(1)

// MARK: - Config

func ConfigFromContext[C any](ctx context.Context) (*C, error) {
	config, err := rootConfigFromContext[C](ctx)
	if err != nil {
		return nil, err
	}

	return &config.Child, nil
}

func rootConfigFromContext[C any](ctx context.Context) (*config[C], error) {
	var c *config[C]

	configValue := ctx.Value(configKey)

	if configValue == nil {
		return c, errors.New("config not found in context")
	}

	if config, ok := configValue.(*config[C]); ok {
		return config, nil
	}

	return c, fmt.Errorf("config has unexpected type: %T", configValue)
}

func contextWithRootConfig[C any](ctx context.Context, config *config[C]) context.Context {
	return context.WithValue(ctx, configKey, config)
}

// MARK: - Logger

var loggerKey = contextKey(0)

var errNoLoggerInContext = errors.New("no logger in context")

func contextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func loggerFromContext(ctx context.Context) (*zap.Logger, error) {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger, nil
	}

	return nil, errNoLoggerInContext
}
