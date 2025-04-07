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

	dotenvParser := dotenv.ParserEnvWithValue(opt.Prefix, ".", transformEnv)

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
	if err := k.Load(env.ProviderWithValue(opt.Prefix, ".", transformEnv), nil); err != nil {
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

func transformEnv(k string, v string) (string, any) {
	// allow specifying nested env vars w/ __
	normalized := strings.ReplaceAll(strings.ToLower(k), "__", ".")
	// split normalized env var by separator
	parts := strings.Split(normalized, ".")
	// pop prefix
	_, parts = parts[0], parts[1:]
	// create final string
	kr := strings.Join(parts, ".")
	kv := transformEnvValue(v)
	return kr, kv
}

func transformEnvValue(v string) any {
	if strings.Contains(v, ",") {
		// if the value contains a comma, split it into a slice
		// each item should be trimmed. if the comma is escaped
		// it should not be split
		escaped := false
		var partsList []string
		var current strings.Builder

		for i := 0; i < len(v); i++ {
			ch := v[i]
			if ch == '\\' && !escaped {
				escaped = true
				continue
			}
			if ch == ',' && !escaped {
				partsList = append(partsList, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteByte(ch)
				escaped = false
			}
		}
		partsList = append(partsList, strings.TrimSpace(current.String()))
		if len(partsList) == 1 {
			return partsList[0]
		}
		return partsList
	}

	return v
}
