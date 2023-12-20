package goji

import (
	"github.com/fruitsco/goji/x/conf"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var envMap = map[string]conf.Environment{
	"development": conf.EnvironmentDevelopment,
	"production":  conf.EnvironmentProduction,
}

func GetEnvFromCLI(ctx *cli.Context) conf.Environment {
	if env, ok := envMap[ctx.String("env")]; ok {
		return env
	}

	return conf.EnvironmentDevelopment
}

func GetLevelFromCLI(ctx *cli.Context) zap.AtomicLevel {
	lvl := ctx.String("log-level")

	if atom, err := zap.ParseAtomicLevel(lvl); err == nil {
		return atom
	}

	env := GetEnvFromCLI(ctx)

	var fallbackLevel zapcore.Level
	if env == conf.EnvironmentDevelopment {
		fallbackLevel = zap.DebugLevel
	} else {
		fallbackLevel = zap.InfoLevel
	}

	return zap.NewAtomicLevelAt(fallbackLevel)
}
