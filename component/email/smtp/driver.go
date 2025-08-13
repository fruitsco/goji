package emailsmtp

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/email"
	"github.com/fruitsco/goji/x/driver"
)

type SMPTDriver struct {
	host string
	port int
}

var _ = email.Driver(&SMPTDriver{})

type SMTPDriverParams struct {
	fx.In

	Config *email.SMTPConfig `optional:"true"`
	Log    *zap.Logger
}

func NewSMTPDriverFactory(params SMTPDriverParams) driver.FactoryResult[email.MailDriver, email.ConnectionFactory] {
	return email.NewConnectionFactory(email.SMTP, func(cfg email.ConnectionConfig) (email.Driver, error) {
		if cfg.Driver != email.SMTP {
			return nil, fmt.Errorf("wrong driver name, expected %s, got %s", email.SMTP, cfg.Driver)
		}

		params.Config = cfg.SMTP

		return NewSMTPDriver(params)
	})
}

// NewSMTPDriver returns a new smtp mailer
func NewSMTPDriver(params SMTPDriverParams) (email.Driver, error) {
	if params.Config == nil {
		return nil, errors.New("config is missing")
	}

	return &SMPTDriver{
		host: params.Config.Host,
		port: params.Config.Port,
	}, nil
}

func (mailer *SMPTDriver) Send(ctx context.Context, msg email.Message) error {
	_, err := mailer.SendID(ctx, msg)
	return err
}

func (mailer *SMPTDriver) SendID(ctx context.Context, msg email.Message) (string, error) {

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
