package goji

import (
	"context"

	"github.com/fruitsco/goji/conf"
	"go.uber.org/zap"
)

type InitParams struct {
	AppName        string
	LogLevel       string
	Prefix         string
	Environment    Environment
	DefaultConfig  conf.DefaultConfig
	ConfigFileName string
}

func Init[C any](ctx context.Context, params InitParams) (context.Context, error) {
	// create the logger
	log, err := createLogger(ctx, params.AppName, params.LogLevel, params.Environment)
	if err != nil {
		return nil, err
	}

	// inject logger into cli context
	ctx = contextWithLogger(ctx, log)

	// parse config using env
	cfg, err := conf.Parse[RootConfig[C]](conf.ParseOptions{
		AppName:     params.AppName,
		Environment: string(params.Environment),
		Defaults:    params.DefaultConfig,
		Prefix:      params.Prefix,
		FileName:    params.ConfigFileName,
		Log:         log,
	})
	if err != nil {
		return nil, err
	}

	// inject the config into the cli context
	ctx = contextWithRootConfig(ctx, cfg)

	return ctx, nil
}

func createLogger(
	ctx context.Context,
	logName string,
	logLevel string,
	environment Environment,
) (*zap.Logger, error) {
	var config zap.Config
	if environment == EnvironmentProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.InitialFields = map[string]any{
		"app": logName,
		"env": environment,
	}

	config.Level = getLogLevel(logLevel, environment)

	return config.Build()
}
