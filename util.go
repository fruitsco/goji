package goji

import (
	"os"

	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var envMap = map[string]Environment{
	"development": EnvironmentDevelopment,
	"production":  EnvironmentProduction,
}

// GetAppEnvFromEnv gets the application environment from the environment variables.
func GetAppEnvFromEnv() Environment {
	v := os.Getenv("ENV")

	if v == "" {
		return EnvironmentDevelopment
	}

	return Environment(v)
}

// GetLogLevelFromEnv gets the log level from the environment variables.
func GetLogLevelFromEnv() string {
	return os.Getenv("LOG_LEVEL")
}

func getEnvFromCLI(cmd *cli.Command) Environment {
	if env, ok := envMap[cmd.String("env")]; ok {
		return env
	}

	return EnvironmentDevelopment
}

func getLogLevel(
	logLevel string,
	environment Environment,
) zap.AtomicLevel {
	if logLevel != "" {
		if atom, err := zap.ParseAtomicLevel(logLevel); err == nil {
			return atom
		}
	}

	var fallbackLevel zapcore.Level
	if environment == EnvironmentProduction {
		fallbackLevel = zap.InfoLevel
	} else {
		fallbackLevel = zap.DebugLevel
	}

	return zap.NewAtomicLevelAt(fallbackLevel)
}
