package email

import (
	"bytes"
	"context"
	"io"

	"github.com/fruitsco/goji/x/driver"
	"github.com/mailgun/mailgun-go/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MailgunDriver struct {
	mg *mailgun.MailgunImpl
}

var _ = Driver(&MailgunDriver{})

type MailgunDriverParams struct {
	fx.In

	Config *MailgunConfig
	Log    *zap.Logger
}

func NewMailgunDriverFactory(params MailgunDriverParams) driver.FactoryResult[MailDriver, Driver] {
	return driver.NewFactory(Mailgun, func() (Driver, error) {
		return NewMailgunDriver(params), nil
	})
}

// NewMailgunDriver returns a new mailgun driver implementation
func NewMailgunDriver(params MailgunDriverParams) *MailgunDriver {
	mg := mailgun.NewMailgun(params.Config.Domain, params.Config.ApiKey)
	mg.SetAPIBase(params.Config.ApiBase)

	return &MailgunDriver{
		mg,
	}
}

// Send a message using mailgun
func (m *MailgunDriver) Send(ctx context.Context, message Message) error {
	text := ""
	if message.Text() != nil {
		text = *message.Text()
	}

	subject := ""
	if message.Subject() != nil {
		subject = *message.Subject()
	}

	from := ""
	if message.From() != nil {
		from = *message.From()
	}

	msg := m.mg.NewMessage(from, subject, text, message.To()...)
	for _, _fl := range message.Files() {
		fl := _fl
		msg.AddReaderInline(fl.Name, io.NopCloser(bytes.NewReader(fl.Data)))
	}

	_, _, err := m.mg.Send(ctx, msg)
	return err
}
