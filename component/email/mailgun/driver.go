package emailmailgun

import (
	"bytes"
	"context"
	"io"

	"github.com/mailgun/mailgun-go/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/email"
	"github.com/fruitsco/goji/x/driver"
)

var Mailgun email.MailDriver = "mailgun"

type MailgunDriver struct {
	mg *mailgun.MailgunImpl
}

var _ = email.Driver(&MailgunDriver{})

type MailgunDriverParams struct {
	fx.In

	Config *email.MailgunConfig
	Log    *zap.Logger
}

func NewMailgunDriverFactory(params MailgunDriverParams) driver.FactoryResult[email.MailDriver, email.Driver] {
	return driver.NewFactory(Mailgun, func() (email.Driver, error) {
		return NewMailgunDriver(params), nil
	})
}

// NewMailgunDriver returns a new mailgun driver implementation
func NewMailgunDriver(params MailgunDriverParams) *MailgunDriver {
	mg := mailgun.NewMailgun(params.Config.Domain, params.Config.ApiKey)
	mg.SetAPIBase(params.Config.ApiBase)

	return &MailgunDriver{mg}
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

	msg := m.mg.NewMessage(from, subject, text, message.GetTo()...)
	for _, _fl := range message.GetFiles() {
		fl := _fl
		msg.AddReaderInline(fl.Name, io.NopCloser(bytes.NewReader(fl.Data)))
	}

	_, id, err := m.mg.Send(ctx, msg)
	return id, err
}
