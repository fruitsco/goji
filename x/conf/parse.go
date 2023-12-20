package conf

import (
	"strings"

	"github.com/fruitsco/go/x/conf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

type ParseOptions struct {
	Environment Environment
	Defaults    conf.DefaultConfig
	Prefix      string
	FileName    string
	Log         *zap.Logger
}

func Parse[C any](opt ParseOptions) (*C, error) {
	var config C

	var log *zap.Logger
	if opt.Log != nil {
		log = opt.Log
	} else {
		log = zap.NewNop()
	}

	k := koanf.New(".")

	k.Load(confmap.Provider(opt.Defaults, "."), nil)

	k.Load(confmap.Provider(conf.DefaultConfig{
		"environment": opt.Environment,
	}, "."), nil)

	if opt.FileName != "" {
		if err := k.Load(file.Provider(opt.FileName), json.Parser()); err != nil {
			log.With(zap.Error(err), zap.String("file", opt.FileName)).Error("error loading file")
		}
	}

	dotenvParser := dotenv.ParserEnv(opt.Prefix, ".", transformEnv)

	if opt.Environment == EnvironmentDevelopment {
		if err := k.Load(file.Provider(".env.dev"), dotenvParser); err != nil {
			log.Debug(".env.dev not found")
		}
	}

	if err := k.Load(file.Provider(".env"), dotenvParser); err != nil {
		log.Debug(".env not found")
	}

	if err := k.Load(env.Provider(opt.Prefix, ".", transformEnv), nil); err != nil {
		log.Error("error loading env vars", zap.Error(err))
		return nil, err
	}

	if err := k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: "conf"}); err != nil {
		log.Error("error unmarshalling config", zap.Error(err))
		return nil, err
	}

	return &config, nil
}

func transformEnv(s string) string {
	// allow specifying nested env vars w/ __
	normalized := strings.ReplaceAll(strings.ToLower(s), "__", ".")
	// split normalized env var by separator
	parts := strings.Split(normalized, ".")
	// pop prefix
	_, parts = parts[0], parts[1:]
	// create final string
	return strings.Join(parts, ".")
}
