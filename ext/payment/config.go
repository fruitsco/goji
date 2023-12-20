package payment

type PaymentDriver string

const (
	Stripe PaymentDriver = "stripe"
)

type StripeConfig struct {
	AccessToken          string  `conf:"access_token"`
	WebhookSecretConnect *string `conf:"webhook_secret_connect"`
	WebhookSecretAccount *string `conf:"webhook_secret_account"`
	InsecureWebhooks     bool    `conf:"insecure_webhooks"`
}

type Config struct {
	Driver PaymentDriver `conf:"driver"`
	Stripe *StripeConfig `conf:"stripe"`
}

var DefaultConfig = map[string]any{
	"payment.driver": "stripe",
}
