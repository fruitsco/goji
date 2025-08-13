package emailresend

import (
	"context"
	"errors"
	"fmt"

	"github.com/resend/resend-go/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/email"
	"github.com/fruitsco/goji/x/driver"
)

type ResendDriver struct {
	r *resend.Client
}

var _ = email.Driver(&ResendDriver{})

type ResendDriverParams struct {
	fx.In

	Config *email.ResendConfig `optional:"true"`
	Log    *zap.Logger
}

// NewResendDriverFactory creates a new resend driver factory
func NewResendDriverFactory(params ResendDriverParams) driver.FactoryResult[email.MailDriver, email.ConnectionFactory] {
	return email.NewConnectionFactory(email.Resend, func(cfg email.ConnectionConfig) (email.Driver, error) {
		if cfg.Driver != email.Mailgun {
			return nil, fmt.Errorf("wrong driver name, expected %s, got %s", email.Resend, cfg.Driver)
		}

		params.Config = cfg.Resend

		return NewResendDriver(params)
	})
}

// NewResendDriver returns a new resend driver implementation
func NewResendDriver(params ResendDriverParams) (email.Driver, error) {
	if params.Config == nil {
		return nil, errors.New("config is missing")
	}

	if params.Config.APIKey == "" {
		return nil, errors.New("api key is empty")
	}

	client := resend.NewClient(params.Config.APIKey)

	return &ResendDriver{client}, nil
}

func (m *ResendDriver) Send(ctx context.Context, message email.Message) error {
	_, err := m.SendID(ctx, message)
	return err
}

func (m *ResendDriver) SendID(ctx context.Context, message email.Message) (string, error) {
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

	var attachments []*resend.Attachment
	for _, _fl := range message.GetFiles() {
		fl := _fl
		attachments = append(attachments, &resend.Attachment{
			Content:  fl.Data,
			Filename: fl.Name,
		})
	}

	req := &resend.SendEmailRequest{
		From:        from,
		To:          message.GetTo(),
		Html:        text,
		Subject:     subject,
		Attachments: attachments,
	}

	// TODO: add idempotency key
	res, err := m.r.Emails.SendWithContext(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}
