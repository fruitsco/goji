package email

import "github.com/fruitsco/goji/conf"

type MailDriver string

const (
	NoOp MailDriver = "noop"
)

type SMTPConfig struct {
	Host string `conf:"host"`
	Port int    `conf:"port"`
}

type MailgunConfig struct {
	Domain  string `conf:"domain"`
	APIKey  string `conf:"api_key"`
	APIBase string `conf:"api_base"`
}

type SenderConfig struct {
	Name *string `conf:"sender_name"`
	Mail *string `conf:"sender_email"`
}

type Config struct {
	Driver  MailDriver     `conf:"driver"`
	Mailgun *MailgunConfig `conf:"mailgun"`
	SMTP    *SMTPConfig    `conf:"smtp"`

	// common config
	Sender *SenderConfig `conf:"sender"`
}

var DefaultConfig = conf.DefaultConfig{
	"email.driver": "noop",
}
