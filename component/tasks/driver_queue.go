package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/queue"
	"github.com/fruitsco/goji/util/randy"
	"github.com/fruitsco/goji/x/driver"
)

type QueueDriver struct {
	queue queue.Queue
	log   *zap.Logger
}

var _ = Driver(&QueueDriver{})

type QueueDriverParams struct {
	fx.In

	Context context.Context
	Queue   queue.Queue
	Log     *zap.Logger
}

func NewQueueDriverFactory(params QueueDriverParams, lc fx.Lifecycle) driver.FactoryResult[TaskDriver, Driver] {
	return driver.NewFactory(Queue, func() (Driver, error) {
		return NewQueueDriver(params, lc)
	})
}

func NewQueueDriver(params QueueDriverParams, lc fx.Lifecycle) (Driver, error) {
	return &QueueDriver{
		queue: params.Queue,
		log:   params.Log.Named("queue"),
	}, nil
}

func (d *QueueDriver) Name() TaskDriver {
	return Queue
}

func (d *QueueDriver) Submit(ctx context.Context, req *CreateTaskRequest) error {
	if req.Queue == "" {
		return fmt.Errorf("queue name / topic is required")
	}

	if req.Data == nil {
		return fmt.Errorf("data is required")
	}

	msg := queue.NewGenericMessage(req.Queue, req.Data)

	for k := range req.Header {
		msg.Meta[k] = req.Header.Get(k)
	}

	name := req.Name
	if name == "" {
		name = fmt.Sprintf("task-%s-%s", req.Queue, randy.Numeric(8))
	}

	scheduleTime := time.Now()
	if req.ScheduleTime != nil {
		scheduleTime = *req.ScheduleTime
	}

	// pubsub does not have the concept of scheduling messages or task names.
	// for sake of consistency, we add these fields to the message metadata.
	// this is not used by the pubsub driver itself, but the information will
	// eventually be delivered to the consumer.
	msg.Meta["task_name"] = name
	msg.Meta["schedule_time"] = scheduleTime.Format(time.RFC3339)

	return d.queue.Publish(ctx, msg)
}

func (d *QueueDriver) Receive(
	ctx context.Context,
	raw RawTask,
) (*Task, error) {
	// the interfaces of queue.PushRequest and tasks.PushRequest match,
	// so we can just cast it here. May diverge in the future.
	message, err := d.queue.Receive(
		ctx,
		queue.NewPushMessageData(raw.GetData(), queue.RawMessageMeta(raw.GetHeader())),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message: %w", err)
	}

	deliveryAttempt := 0
	if message.GetDeliveryAttempt() != nil {
		deliveryAttempt = *message.GetDeliveryAttempt()
	}

	// pubsub does not have the concept of scheduling messages. this is actually
	// the reason why we implement cloud tasks in the first place. so when using
	// pubsub as a fallback (e.g. for local development), we just set the schedule
	// time to the publish time.
	scheduleTime := message.GetPublishTime()

	// extract task name from message metadata. this is not used by pubsub itself,
	// but rather attached when publishing for consistency with other drivers.
	meta := message.GetMeta()

	queueName := meta["queue_name"]
	taskName := meta["task_name"]

	header := make(http.Header)
	for k, v := range message.GetMeta() {
		header.Set(k, v)
	}

	return &Task{
		QueueName:      queueName,
		TaskName:       taskName,
		ScheduleTime:   scheduleTime,
		RetryCount:     deliveryAttempt,
		ExecutionCount: deliveryAttempt,
		Data:           message.GetData(),
		Header:         header,
	}, nil
}
