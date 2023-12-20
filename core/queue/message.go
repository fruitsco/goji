package queue

import "time"

type Message interface {
	GetTopic() string
	GetData() []byte
	GetID() string
	GetPublishTime() time.Time
	GetDeliveryAttempt() *int
}

type GenericMessage struct {
	Message

	ID              string
	Topic           string
	Data            []byte
	PublishTime     time.Time
	DeliveryAttempt *int
}

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

func NewGenericMessage(topic string, data []byte) *GenericMessage {
	return &GenericMessage{
		Topic: topic,
		Data:  data,
	}
}
