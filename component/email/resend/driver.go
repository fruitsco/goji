package emailresend

import (
	"context"

	"github.com/resend/resend-go/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/email"
	"github.com/fruitsco/goji/x/driver"
)

// NewResend creates a new resend client
func NewResend(config *email.ResendConfig) *resend.Client {
	return resend.NewClient(config.APIKey)
}

type ResendDriver struct {
	r *resend.Client
}

var _ = email.Driver(&ResendDriver{})

type ResendDriverParams struct {
	fx.In

	Resend *resend.Client
	Log    *zap.Logger
}

// NewResendDriverFactory creates a new resend driver factory
func NewResendDriverFactory(params ResendDriverParams) driver.FactoryResult[email.MailDriver, email.Driver] {
	return driver.NewFactory(email.Resend, func() (email.Driver, error) {
		return NewResendDriver(params), nil
	})
}

// NewResendDriver returns a new resend driver implementation
func NewResendDriver(params ResendDriverParams) *ResendDriver {
	return &ResendDriver{params.Resend}
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
