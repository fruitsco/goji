package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SlackDriver struct {
	config *SlackConfig
	log    *zap.Logger
}

var _ = Driver(&SlackDriver{})

type SlackDriverParams struct {
	fx.In

	Config *SlackConfig
	Log    *zap.Logger
}

func NewSlackDriverFactory(params SlackDriverParams) driver.FactoryResult[NotificationDriver, Driver] {
	return driver.NewFactory(Slack, func() (Driver, error) {
		return NewSlackDriver(params), nil
	})
}

// NewPostinoDriver returns a new queue mailer implementation
func NewSlackDriver(params SlackDriverParams) *SlackDriver {
	return &SlackDriver{
		config: params.Config,
		log:    params.Log,
	}
}

type Payload struct {
	Text        string    `json:"text,omitempty"`
	Attachments []Message `json:"attachments,omitempty"`
}

func (m *SlackDriver) Send(ctx context.Context, message Message) error {
	// Convert the payload to JSON
	payload := Payload{
		Text: message.Text,
		Attachments: []Message{
			message,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", m.config.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		m.log.Error("Error sending Domain Upload Notifcation", zap.String("status", resp.Status))
		return errors.New(fmt.Sprintf("Error sending Domain Upload Notifcation. Status: %v", resp.Status))
	}
	log.Printf("Response Status: %v", resp.Status)

	return nil
}
