package tasks

import (
	"context"
	"fmt"
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

func NewQueueDriverFactory(params CloudTasksDriverParams, lc fx.Lifecycle) driver.FactoryResult[TaskDriver, Driver] {
	return driver.NewFactory(CloudTasks, func() (Driver, error) {
		return NewCloudTasksDriver(params, lc)
	})
}

func NewQueueDriver(params QueueDriverParams, lc fx.Lifecycle) (Driver, error) {
	// NOTE: Cloud Tasks does not have an emulator (yet)
	// if params.Config != nil && params.Config.EmulatorHost != nil {
	// 	os.Setenv("PUBSUB_EMULATOR_HOST", *params.Config.EmulatorHost)
	// 	os.Setenv("PUBSUB_PROJECT_ID", params.Config.ProjectID)
	// }

	return &QueueDriver{
		queue: params.Queue,
		log:   params.Log.Named("queue"),
	}, nil
}

func (d *QueueDriver) Name() TaskDriver {
	return Queue
}

func (d *QueueDriver) Submit(ctx context.Context, req CreateTaskRequest) error {
	queueReq, ok := req.(*CreateQueueTaskRequest)
	if !ok {
		return fmt.Errorf("invalid request type, expected *CreateHttpTaskRequest, got %T", req)
	}

	if queueReq.Topic == "" {
		return fmt.Errorf("topic is required")
	}

	if queueReq.Data == nil {
		return fmt.Errorf("data is required")
	}

	message := queue.NewGenericMessage(queueReq.Topic, queueReq.Data)

	name := queueReq.Name
	if name == "" {
		name = fmt.Sprint("task-", randy.Numeric(8))
	}

	scheduleTime := time.Now()
	if queueReq.ScheduleTime != nil {
		scheduleTime = *queueReq.ScheduleTime
	}

	// pubsub does not have the concept of scheduling messages or task names.
	// for sake of consistency, we add these fields to the message metadata.
	// this is not used by the pubsub driver itself, but the information will
	// eventually be delivered to the consumer.
	message.Meta["queue_name"] = queueReq.Topic
	message.Meta["task_name"] = name
	message.Meta["schedule_time"] = scheduleTime.Format(time.RFC3339)

	return d.queue.Publish(ctx, message)
}

func (d *QueueDriver) ReceivePush(
	ctx context.Context,
	req PushRequest,
) (*Task, error) {
	// the interfaces of queue.PushRequest and tasks.PushRequest match,
	// so we can just cast it here. May diverge in the future.
	message, err := d.queue.ReceivePush(ctx, queue.PushRequest{
		Data:    req.Data,
		Headers: req.Meta,
	})
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

	return &Task{
		QueueName:      queueName,
		TaskName:       taskName,
		ScheduleTime:   scheduleTime,
		RetryCount:     deliveryAttempt,
		ExecutionCount: deliveryAttempt,
		Data:           message.GetData(),
		Meta:           message.GetMeta(),
	}, nil
}
