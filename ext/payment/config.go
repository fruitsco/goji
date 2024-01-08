package payment

import "github.com/fruitsco/goji/x/conf"

type PaymentDriver string

const (
	Stripe PaymentDriver = "stripe"
)

type Config struct {
	Driver PaymentDriver `conf:"driver"`
	Stripe *StripeConfig `conf:"stripe"`
}

var DefaultConfig = conf.DefaultConfig{
	"payment.driver": "stripe",
}
