package goji

import (
	"context"

	"github.com/fruitsco/goji/conf"
)

type InitParams struct {
	AppName        string
	Prefix         string
	Environment    Environment
	DefaultConfig  conf.DefaultConfig
	ConfigFileName string
}

func Init(ctx context.Context, params AppParams) context.Context {
	// create the logger
	log, err := createLogger(ctx, params.AppName, params.Environment)
	if err != nil {
		return err
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
		return err
	}

	// inject the config into the cli context
	ctx = contextWithRootConfig(ctx, cfg)

	return ctx
}
