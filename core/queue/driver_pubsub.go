package queue

import (
	"context"
	"encoding/json"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/fruitsco/goji/x/driver"
	"github.com/fruitsco/roma/x/google"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type pubSubPushMessage struct {
	Subscription string         `json:"subscription"`
	Message      pubsub.Message `json:"message"`
}

type PubSubDriver struct {
	client *pubsub.Client
	log    *zap.Logger
}

var _ = Driver(&PubSubDriver{})

type PubSubDriverParams struct {
	fx.In

	Context context.Context
	Config  *PubSubConfig
	Log     *zap.Logger
}

func NewPubSubDriverFactory(params PubSubDriverParams) driver.FactoryResult[QueueDriver, Driver] {
	return driver.NewFactory(PubSub, func() (Driver, error) {
		return NewPubSubDriver(params)
	})
}

func NewPubSubDriver(params PubSubDriverParams) (*PubSubDriver, error) {
	credentials := google.NewCredentials(params.Context, []string{
		pubsub.ScopePubSub,
		pubsub.ScopeCloudPlatform,
	})

	if params.Config != nil && params.Config.EmulatorHost != nil {
		// put the configuration into env. pubsub client expects them there.
		os.Setenv("PUBSUB_EMULATOR_HOST", *params.Config.EmulatorHost)
		os.Setenv("PUBSUB_PROJECT_ID", params.Config.ProjectID)
	}

	opts := make([]option.ClientOption, 0)

	if clientOption := credentials.ClientOption(); clientOption != nil {
		opts = append(opts, clientOption)
	}

	client, err := pubsub.NewClient(
		params.Context,
		params.Config.ProjectID,
		opts...,
	)

	if err != nil {
		return nil, err
	}

	return &PubSubDriver{
		client: client,
		log:    params.Log.Named("pubsub"),
	}, nil
}

func (q *PubSubDriver) Name() QueueDriver {
	return PubSub
}

func (q *PubSubDriver) Publish(ctx context.Context, message Message) error {
	topic := q.client.Topic(message.GetTopic())

	q.log.With(zap.String("topic", message.GetTopic())).Debug("publishing message")
	_, err := topic.Publish(ctx, &pubsub.Message{
		Data: message.GetData(),
	}).Get(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (q *PubSubDriver) RecievePush(ctx context.Context, topic string, data []byte) (Message, error) {
	message := &pubSubPushMessage{}
	err := json.Unmarshal(data, message)

	if err != nil {
		return nil, err
	}

	return &GenericMessage{
		ID:              message.Message.ID,
		Data:            message.Message.Data,
		DeliveryAttempt: message.Message.DeliveryAttempt,
		PublishTime:     message.Message.PublishTime,
		Topic:           topic,
	}, nil
}
