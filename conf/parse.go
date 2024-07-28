package conf

import (
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	jsonParser "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type ParseOptions struct {
	AppName     string
	Environment string
	Defaults    DefaultConfig
	Prefix      string
	Log         *zap.Logger
	FileName    string
}

func Parse[C any](opt ParseOptions) (*C, error) {
	var config C

	log := zap.L()
	if opt.Log != nil {
		log = opt.Log
	}

	k := koanf.New(".")

	dotenvParser := dotenv.ParserEnv(opt.Prefix, ".", transformEnv)

	// PRIO 0 - defaults
	k.Load(confmap.Provider(opt.Defaults, "."), nil)

	// PRIO 1 - config file
	if opt.FileName != "" {
		if err := k.Load(file.Provider(opt.FileName), jsonParser.Parser()); err != nil {
			log.With(zap.Error(err), zap.String("file", opt.FileName)).Error("error loading file")
		}
	}

	// PRIO 2 - load .env.dev if in development environment
	if opt.Environment != "production" {
		if err := k.Load(file.Provider(".env.dev"), dotenvParser); err != nil {
			log.Debug(".env.dev not found")
		}
	}

	// PRIO 3 - load .env.test if in test environment
	if opt.Environment == "test" {
		if err := k.Load(file.Provider(".env.test"), dotenvParser); err != nil {
			log.Debug(".env.test not found")
		}
	}

	// PRIO 4 - load .env
	if err := k.Load(file.Provider(".env"), dotenvParser); err != nil {
		log.Debug(".env not found")
	}

	// PRIO 5 - load env vars
	if err := k.Load(env.Provider(opt.Prefix, ".", transformEnv), nil); err != nil {
		log.Error("error loading env vars", zap.Error(err))
		return nil, err
	}

	// PRIO 6 - set actual environment
	k.Load(confmap.Provider(DefaultConfig{
		"app.name": opt.AppName,
		"app.env":  opt.Environment,
	}, "."), nil)

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
