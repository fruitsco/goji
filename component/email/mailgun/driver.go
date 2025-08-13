package emailmailgun

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/mailgun/mailgun-go/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/email"
	"github.com/fruitsco/goji/x/driver"
)

// NewMailgun creates a new mailgun client
func NewMailgun(config *email.MailgunConfig) mailgun.Mailgun {
	mg := mailgun.NewMailgun(config.APIKey)
	mg.SetAPIBase(config.APIBase)
	return mg
}

type MailgunDriver struct {
	mg     mailgun.Mailgun
	domain string
}

var _ = email.Driver(&MailgunDriver{})

type MailgunDriverParams struct {
	fx.In

	Config *email.MailgunConfig `optional:"true"`
	Log    *zap.Logger
}

// NewMailgunDriverFactory creates a new mailgun driver factory
func NewMailgunDriverFactory(params MailgunDriverParams) driver.FactoryResult[email.MailDriver, email.ConnectionFactory] {
	return email.NewConnectionFactory(email.Mailgun, func(cfg email.ConnectionConfig) (email.Driver, error) {
		params.Config = cfg.Mailgun

		return NewMailgunDriver(params)
	})
}

// NewMailgunDriver returns a new mailgun driver implementation
func NewMailgunDriver(params MailgunDriverParams) (*MailgunDriver, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("config is missing")
	}

	if params.Config.APIKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}

	if params.Config.APIBase == "" {
		return nil, fmt.Errorf("api base is empty")
	}

	mg := NewMailgun(params.Config)

	return &MailgunDriver{mg, params.Config.Domain}, nil
}

func (m *MailgunDriver) Send(ctx context.Context, message email.Message) error {
	_, err := m.SendID(ctx, message)
	return err
}

func (m *MailgunDriver) SendID(ctx context.Context, message email.Message) (string, error) {
	text := ""
	if message.GetText() != nil {
		text = *message.GetText()
	}

	subject := ""
	if message.GetSubject() != nil {
		subject = *message.GetSubject()
	}

	from := ""
	if message.GetFrom() != nil {
		from = *message.GetFrom()
	}

	msg := mailgun.NewMessage(m.domain, from, subject, text, message.GetTo()...)
	for _, _fl := range message.GetFiles() {
		fl := _fl
		msg.AddReaderInline(fl.Name, io.NopCloser(bytes.NewReader(fl.Data)))
	}

	res, err := m.mg.Send(ctx, msg)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}
