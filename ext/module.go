package core

import (
	"github.com/fruitsco/goji/ext/payment"
	"go.uber.org/fx"
)

func Module(config *Config) fx.Option {
	return fx.Module("ext",
		payment.Module(config.Payment),
	)
}
