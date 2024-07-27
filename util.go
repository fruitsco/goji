package goji

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var envMap = map[string]Environment{
	"development": EnvironmentDevelopment,
	"production":  EnvironmentProduction,
}

func getEnvFromCLI(ctx *cli.Context) Environment {
	if env, ok := envMap[ctx.String("env")]; ok {
		return env
	}

	return EnvironmentDevelopment
}

func getLevelFromCLI(ctx *cli.Context) zap.AtomicLevel {
	lvl := ctx.String("log-level")

	if lvl != "" {
		if atom, err := zap.ParseAtomicLevel(lvl); err == nil {
			return atom
		}
	}

	env := getEnvFromCLI(ctx)

	var fallbackLevel zapcore.Level
	if env == EnvironmentProduction {
		fallbackLevel = zap.InfoLevel
	} else {
		fallbackLevel = zap.DebugLevel
	}

	return zap.NewAtomicLevelAt(fallbackLevel)
}
