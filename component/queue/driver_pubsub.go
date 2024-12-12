package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/fruitsco/goji/internal/google"
	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// pubSubPushMessage is a struct that represents the message that is sent to the
// push endpoint of a pubsub subscription. It contains the subscription name and
// the message itself. See https://cloud.google.com/pubsub/docs/push#receive_push
type pubSubPushMessage struct {
	Subscription string         `json:"subscription"`
	Message      pubsub.Message `json:"message"`
}

type PubSubDriver struct {
	topicMap map[string]*pubsub.Topic
	client   *pubsub.Client
	log      *zap.Logger
}

var _ = Driver(&PubSubDriver{})

type PubSubDriverParams struct {
	fx.In

	Context context.Context
	Config  *PubSubConfig
	Log     *zap.Logger
}

func NewPubSubDriverFactory(params PubSubDriverParams, lc fx.Lifecycle) driver.FactoryResult[QueueDriver, Driver] {
	return driver.NewFactory(PubSub, func() (Driver, error) {
		return NewPubSubDriver(params, lc)
	})
}

func NewPubSubDriver(params PubSubDriverParams, lc fx.Lifecycle) (*PubSubDriver, error) {
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

	driver := &PubSubDriver{
		client:   client,
		topicMap: make(map[string]*pubsub.Topic),
		log:      params.Log.Named("pubsub"),
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return driver.Close()
		},
	})

	return driver, nil
}

func (q *PubSubDriver) Name() QueueDriver {
	return PubSub
}

func (q *PubSubDriver) Publish(ctx context.Context, message Message) error {
	topic := q.getTopic(message.GetTopic())

	_, err := topic.Publish(ctx, &pubsub.Message{Data: message.GetData()}).Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (q *PubSubDriver) Receive(ctx context.Context, raw RawMessage) (Message, error) {
	message := &pubSubPushMessage{}

	if err := json.Unmarshal(raw.GetData(), message); err != nil {
		return nil, err
	}

	return &GenericMessage{
		ID:              message.Message.ID,
		Data:            message.Message.Data,
		DeliveryAttempt: message.Message.DeliveryAttempt,
		PublishTime:     message.Message.PublishTime,
		Meta:            message.Message.Attributes,
	}, nil
}

func (q *PubSubDriver) Close() error {
	// first, stop all open topics
	for _, topic := range q.topicMap {
		topic.Stop()
	}

	// then close the client
	return q.client.Close()
}

// The pubsub client's Topic() method's documentation states:
// "Avoid creating many Topic instances if you use them to publish."
// This is why we are caching the topic instances.
func (q *PubSubDriver) getTopic(topic string) *pubsub.Topic {
	if t, ok := q.topicMap[topic]; ok {
		return t
	}

	t := q.client.Topic(topic)
	q.topicMap[topic] = t

	return t
}
