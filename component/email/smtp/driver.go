package emailsmtp

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/email"
	"github.com/fruitsco/goji/x/driver"
)

var Smtp email.MailDriver = "smtp"

type SmtpDriver struct {
	host string
	port int
}

var _ = email.Driver(&SmtpDriver{})

type SmtpDriverParams struct {
	fx.In

	Config *email.SmtpConfig
	Log    *zap.Logger
}

func NewSmtpDriverFactory(params SmtpDriverParams) driver.FactoryResult[email.MailDriver, email.Driver] {
	return driver.NewFactory(Smtp, func() (email.Driver, error) {
		return NewSmtpDriver(params), nil
	})
}

// NewSmtpDriver returns a new smtp mailer
func NewSmtpDriver(params SmtpDriverParams) *SmtpDriver {
	return &SmtpDriver{
		host: params.Config.Host,
		port: params.Config.Port,
	}
}

func (mailer *SmtpDriver) Send(ctx context.Context, msg email.Message) error {
	_, err := mailer.SendID(ctx, msg)
	return err
}

func (mailer *SmtpDriver) SendID(ctx context.Context, msg email.Message) (string, error) {

	smtpAddr := fmt.Sprintf("%s:%d", mailer.host, mailer.port)

	text := ""
	if msg.GetText() != nil {
		text = *msg.GetText()
	}

	if err := smtp.SendMail(smtpAddr, nil, *msg.GetFrom(), msg.GetTo(), []byte(text)); err != nil {
		fmt.Println(err)
		return "", err
	}

	return uuid.NewString(), nil
}
