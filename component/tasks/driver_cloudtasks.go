package tasks

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
)

type CloudTasksDriver struct {
	config *CloudTasksConfig
	client *cloudtasks.Client
	log    *zap.Logger
}

var _ = Driver(&CloudTasksDriver{})

type CloudTasksDriverParams struct {
	fx.In

	Context context.Context
	Config  *CloudTasksConfig
	Log     *zap.Logger
}

func NewCloudTasksDriverFactory(params CloudTasksDriverParams, lc fx.Lifecycle) driver.FactoryResult[TaskDriver, Driver] {
	return driver.NewFactory(CloudTasks, func() (Driver, error) {
		return NewCloudTasksDriver(params, lc)
	})
}

func NewCloudTasksDriver(params CloudTasksDriverParams, lc fx.Lifecycle) (Driver, error) {
	// NOTE: Cloud Tasks does not have an emulator (yet)
	// if params.Config != nil && params.Config.EmulatorHost != nil {
	// 	os.Setenv("PUBSUB_EMULATOR_HOST", *params.Config.EmulatorHost)
	// 	os.Setenv("PUBSUB_PROJECT_ID", params.Config.ProjectID)
	// }

	client, err := cloudtasks.NewClient(params.Context)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return &CloudTasksDriver{
		client: client,
		config: params.Config,
		log:    params.Log.Named("cloudtasks"),
	}, nil
}

func (d *CloudTasksDriver) Name() TaskDriver {
	return CloudTasks
}

func (d *CloudTasksDriver) Submit(ctx context.Context, req CreateTaskRequest) error {
	httpReq, ok := req.(*CreateHttpTaskRequest)
	if !ok {
		return fmt.Errorf("invalid request type, expected *CreateHttpTaskRequest, got %T", req)
	}

	if httpReq.Queue == "" {
		return fmt.Errorf("queue name is required")
	}

	if httpReq.Url == "" {
		return fmt.Errorf("url is required")
	}

	if httpReq.Method == "" {
		httpReq.Method = http.MethodPost
	}

	if _, ok := httpMethodMap[httpReq.Method]; !ok {
		return fmt.Errorf("invalid http method: %s", httpReq.Method)
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", d.config.ProjectID, d.config.Region, httpReq.Queue)

	taskReq := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: httpMethodMap[httpReq.Method],
					Url:        httpReq.Url,
					Body:       httpReq.Body,
					Headers:    httpReq.Headers,
				},
			},
		},
	}

	if _, err := d.client.CreateTask(ctx, taskReq); err != nil {
		return fmt.Errorf("could not create task: %v", err)
	}

	return nil
}

func (d *CloudTasksDriver) ReceivePush(
	ctx context.Context,
	req PushRequest,
) (*Task, error) {
	scheduleTimeValue, ok := req.Meta["X-CloudTasks-TaskETA"]
	if !ok {
		return nil, fmt.Errorf("missing schedule time")
	}

	scheduleTimeSeconds, err := strconv.ParseInt(scheduleTimeValue, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule time: %v", err)
	}

	retryCount := 0
	if retryCountValue, ok := req.Meta["X-CloudTasks-TaskRetryCount"]; ok {
		retryCount, err = strconv.Atoi(retryCountValue)
		if err != nil {
			return nil, fmt.Errorf("invalid retry count: %v", err)
		}
	}

	executionCount := 0
	if executionCountValue, ok := req.Meta["X-CloudTasks-TaskExecutionCount"]; ok {
		executionCount, err = strconv.Atoi(executionCountValue)
		if err != nil {
			return nil, fmt.Errorf("invalid execution count: %v", err)
		}
	}

	return &Task{
		TaskName:       req.TaskName,
		Data:           req.Data,
		ScheduleTime:   time.Unix(scheduleTimeSeconds, 0),
		RetryCount:     retryCount,
		ExecutionCount: executionCount,
		Meta:           req.Meta,
	}, nil
}

// Header	Description
// X-CloudTasks-QueueName	The name of the queue.
// X-CloudTasks-TaskName	The "short" name of the task, or, if no name was specified at creation, a unique system-generated id. This is the my-task-id value in the complete task name, ie, task_name = projects/my-project-id/locations/my-location/queues/my-queue-id/tasks/my-task-id.
// X-CloudTasks-TaskRetryCount	The number of times this task has been retried. For the first attempt, this value is 0. This number includes attempts where the task failed due to 5XX error codes and never reached the execution phase.
// X-CloudTasks-TaskExecutionCount	The total number of times that the task has received a response from the handler. Since Cloud Tasks deletes the task once a successful response has been received, all previous handler responses were failures. This number does not include failures due to 5XX error codes.
// X-CloudTasks-TaskETA	The schedule time of the task, specified in seconds since January 1st 1970.
// In addition, requests from Cloud Tasks might contain the following headers:

// Header	Description
// X-CloudTasks-TaskPreviousResponse	The HTTP response code from the previous retry.
// X-CloudTasks-TaskRetryReason	The reason for retrying the task.

var httpMethodMap = map[string]taskspb.HttpMethod{
	http.MethodGet:     taskspb.HttpMethod_GET,
	http.MethodPost:    taskspb.HttpMethod_POST,
	http.MethodPut:     taskspb.HttpMethod_PUT,
	http.MethodPatch:   taskspb.HttpMethod_PATCH,
	http.MethodDelete:  taskspb.HttpMethod_DELETE,
	http.MethodHead:    taskspb.HttpMethod_HEAD,
	http.MethodOptions: taskspb.HttpMethod_OPTIONS,
}
