package email

type MailDriver string

const (
	Mailgun MailDriver = "mailgun"
	Smtp    MailDriver = "smtp"
	NoOp    MailDriver = "noop"
)

type SmtpConfig struct {
	Host string `conf:"host"`
	Port int    `conf:"port"`
}

type MailgunConfig struct {
	Domain  string `conf:"domain"`
	ApiKey  string `conf:"api_key"`
	ApiBase string `conf:"api_base"`
}

type SenderConfig struct {
	Name *string `conf:"sender_name"`
	Mail *string `conf:"sender_email"`
}

type Config struct {
	Driver  MailDriver     `conf:"driver"`
	Mailgun *MailgunConfig `conf:"mailgun"`
	Smtp    *SmtpConfig    `conf:"smtp"`

	// common config
	Sender *SenderConfig `conf:"sender"`
}

var DefaultConfig = map[string]any{
	"email.driver": "noop",
}
