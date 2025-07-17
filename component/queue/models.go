package queue

import (
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
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
		ID:    uuid.NewString(),
		Topic: topic,
		Data:  data,
		Meta:  make(map[string]string),
	}
}

type RawMessageMeta map[string][]string

type RawMessage interface {
	GetData() []byte
	GetMeta() RawMessageMeta
}

type PushMessageData struct {
	data []byte
	meta RawMessageMeta
}

func (m *PushMessageData) GetData() []byte {
	return m.data
}

func (m *PushMessageData) GetMeta() RawMessageMeta {
	return m.meta
}

func NewPushMessageData(data []byte, meta RawMessageMeta) *PushMessageData {
	return &PushMessageData{
		data: data,
		meta: meta,
	}
}

func NewPushMessageDataFromRequest(r *http.Request) (*PushMessageData, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	meta := make(map[string][]string)
	for k, v := range r.Header {
		meta[k] = v
	}

	return NewPushMessageData(data, meta), nil
}
