package core

import (
	"github.com/fruitsco/goji/ext/payment"
	"github.com/fruitsco/goji/util"
)

type Config struct {
	Payment *payment.Config `conf:"payment"`
}

func NewConfig(
	payment *payment.Config,
) Config {
	return Config{
		Payment: payment,
	}
}

var DefaultConfig = util.MergeMap(
	payment.DefaultConfig,
)
