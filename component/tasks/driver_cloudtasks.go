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
	"google.golang.org/protobuf/types/known/timestamppb"

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

	if params.Config == nil {
		return nil, fmt.Errorf("missing cloud tasks config")
	}

	if params.Config.ProjectID == "" {
		return nil, fmt.Errorf("cloudtasks is missing project id")
	}

	if params.Config.Region == "" {
		return nil, fmt.Errorf("cloudtasks is missing region")
	}

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

	var taskName string
	if req.GetName() != "" {
		taskName = fmt.Sprintf("%s/tasks/%s", queuePath, req.GetName())
	}

	var scheduleTime *timestamppb.Timestamp
	if req.GetScheduleTime() != nil {
		scheduleTime = timestamppb.New(*req.GetScheduleTime())
	}

	var authHeader *taskspb.HttpRequest_OidcToken
	if d.config.AuthServiceAccountEmail != "" {
		authHeader = &taskspb.HttpRequest_OidcToken{
			OidcToken: &taskspb.OidcToken{
				ServiceAccountEmail: d.config.AuthServiceAccountEmail,
			},
		}
	}

	taskReq := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			Name:         taskName,
			ScheduleTime: scheduleTime,
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod:          httpMethodMap[httpReq.Method],
					Url:                 httpReq.Url,
					Body:                httpReq.Body,
					Headers:             httpReq.Headers,
					AuthorizationHeader: authHeader,
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
	taskName, ok := req.Meta["X-CloudTasks-TaskName"]
	if !ok {
		return nil, fmt.Errorf("invalid cloud tasks request: missing task name")
	}

	queueName, ok := req.Meta["X-CloudTasks-QueueName"]
	if !ok {
		return nil, fmt.Errorf("invalid cloud tasks request: missing queue name")
	}

	scheduleTimeValue, ok := req.Meta["X-CloudTasks-TaskETA"]
	if !ok {
		return nil, fmt.Errorf("invalid cloud tasks request: missing schedule time")
	}

	scheduleTimeSeconds, err := strconv.ParseInt(scheduleTimeValue, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid cloud tasks request: invalid schedule time: %v", err)
	}

	retryCount := 0
	if retryCountValue, ok := req.Meta["X-CloudTasks-TaskRetryCount"]; ok {
		retryCount, err = strconv.Atoi(retryCountValue)
		if err != nil {
			return nil, fmt.Errorf("invalid cloud tasks request: invalid retry count: %v", err)
		}
	}

	executionCount := 0
	if executionCountValue, ok := req.Meta["X-CloudTasks-TaskExecutionCount"]; ok {
		executionCount, err = strconv.Atoi(executionCountValue)
		if err != nil {
			return nil, fmt.Errorf("invalid cloud tasks request: invalid execution count: %v", err)
		}
	}

	return &Task{
		TaskName:       taskName,
		QueueName:      queueName,
		Data:           req.Data,
		ScheduleTime:   time.Unix(scheduleTimeSeconds, 0),
		RetryCount:     retryCount,
		ExecutionCount: executionCount,
		Meta:           req.Meta,
	}, nil
}

var httpMethodMap = map[string]taskspb.HttpMethod{
	http.MethodGet:     taskspb.HttpMethod_GET,
	http.MethodPost:    taskspb.HttpMethod_POST,
	http.MethodPut:     taskspb.HttpMethod_PUT,
	http.MethodPatch:   taskspb.HttpMethod_PATCH,
	http.MethodDelete:  taskspb.HttpMethod_DELETE,
	http.MethodHead:    taskspb.HttpMethod_HEAD,
	http.MethodOptions: taskspb.HttpMethod_OPTIONS,
}
