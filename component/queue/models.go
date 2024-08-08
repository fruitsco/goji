package queue

import (
	"net/http"
	"time"
)

type Message interface {
	GetTopic() string
	GetData() []byte
	GetID() string
	GetPublishTime() time.Time
	GetDeliveryAttempt() *int
	GetMeta() map[string]string
}

type GenericMessage struct {
	ID              string
	Topic           string
	Data            []byte
	PublishTime     time.Time
	DeliveryAttempt *int
	Meta            map[string]string
}

var _ = Message(&GenericMessage{})

func (m *GenericMessage) GetTopic() string {
	return m.Topic
}

func (m *GenericMessage) GetData() []byte {
	return m.Data
}

func (m *GenericMessage) GetID() string {
	return m.ID
}

func (m *GenericMessage) GetPublishTime() time.Time {
	return m.PublishTime
}

func (m *GenericMessage) GetDeliveryAttempt() *int {
	return m.DeliveryAttempt
}

func (m *GenericMessage) GetMeta() map[string]string {
	return m.Meta
}

func NewGenericMessage(topic string, data []byte) *GenericMessage {
	return &GenericMessage{
		Topic: topic,
		Data:  data,
		Meta:  make(map[string]string),
	}
}

type PushRequest struct {
	Data   []byte
	Header http.Header
}
