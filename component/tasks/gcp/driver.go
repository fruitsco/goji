package tasksgcp

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fruitsco/goji"
	"github.com/fruitsco/goji/component/tasks"
	"github.com/fruitsco/goji/x/driver"
)

type CloudTasksDriver struct {
	config *tasks.CloudTasksConfig
	client *cloudtasks.Client
	log    *zap.Logger
}

var _ = tasks.Driver(&CloudTasksDriver{})

type CloudTasksDriverParams struct {
	fx.In

	// Context is the context to use for the driver.
	Context context.Context

	// Config is the cloud tasks configuration.
	Config *tasks.CloudTasksConfig

	// GRPCConn is the gRPC connection to use for the driver.
	GRPCConn *grpc.ClientConn `optional:"true"`

	// NoAuth is a flag to disable authentication.
	// This flag should be set to `true` only for testing purposes.
	NoAuth bool `optional:"true"`

	// Log is the logger to use for the driver.
	Log *zap.Logger
}

func NewCloudTasksDriverFactory(params CloudTasksDriverParams, lc fx.Lifecycle) driver.FactoryResult[tasks.TaskDriver, tasks.Driver] {
	factory := driver.NewFactory(tasks.CloudTasks, func() (tasks.Driver, error) {
		driver, err := NewCloudTasksDriver(params)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				return driver.Close()
			},
		})

		return driver, nil
	})

	return factory
}

func NewCloudTasksDriver(params CloudTasksDriverParams) (*CloudTasksDriver, error) {
	// NOTE: Cloud Tasks does not have an emulator (yet)
	// if params.Config != nil && params.Config.EmulatorHost != nil {
	// 	os.Setenv("PUBSUB_EMULATOR_HOST", *params.Config.EmulatorHost)
	// 	os.Setenv("PUBSUB_PROJECT_ID", params.Config.ProjectID)
	// }

	if params.Config == nil {
		return nil, fmt.Errorf("missing cloudtasks config")
	}

	if params.Config.ProjectID == "" {
		return nil, fmt.Errorf("cloudtasks is missing project id")
	}

	if params.Config.Region == "" {
		return nil, fmt.Errorf("cloudtasks is missing region")
	}

	options := make([]option.ClientOption, 0)

	if params.NoAuth {
		options = append(options, option.WithoutAuthentication())
	}

	if params.GRPCConn != nil {
		options = append(options, option.WithGRPCConn(params.GRPCConn))
	}

	client, err := cloudtasks.NewClient(params.Context, options...)
	if err != nil {
		return nil, err
	}

	return &CloudTasksDriver{
		client: client,
		config: params.Config,
		log:    params.Log.Named("cloudtasks"),
	}, nil
}

func (d *CloudTasksDriver) Name() tasks.TaskDriver {
	return tasks.CloudTasks
}

func (d *CloudTasksDriver) Submit(ctx context.Context, req *tasks.CreateTaskRequest) error {
	if req.Queue == "" {
		return fmt.Errorf("queue name is required")
	}

	url := d.config.DefaultUrl
	if req.Url != "" {
		url = req.Url
	}
	if url == "" {
		return fmt.Errorf("url is required")
	}

	if req.Method == "" {
		req.Method = http.MethodPost
	}
	if _, ok := httpMethodMap[req.Method]; !ok {
		return fmt.Errorf("invalid http method: %s", req.Method)
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", d.config.ProjectID, d.config.Region, req.Queue)

	var taskName string
	if req.Name != "" {
		taskName = fmt.Sprintf("%s/tasks/%s", queuePath, req.Name)
	}

	var scheduleTime *timestamppb.Timestamp
	if req.ScheduleTime != nil {
		scheduleTime = timestamppb.New(*req.ScheduleTime)
	}

	var authHeader *taskspb.HttpRequest_OidcToken
	if d.config.AuthServiceAccountEmail != "" {
		authHeader = &taskspb.HttpRequest_OidcToken{
			OidcToken: &taskspb.OidcToken{
				ServiceAccountEmail: d.config.AuthServiceAccountEmail,
			},
		}
	}

	headers := make(map[string]string)
	for k, v := range req.Header {
		headers[k] = strings.Join(v, ",")
	}

	taskReq := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			Name:         taskName,
			ScheduleTime: scheduleTime,
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod:          httpMethodMap[req.Method],
					Url:                 url,
					Body:                req.Data,
					Headers:             headers,
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

func (d *CloudTasksDriver) Receive(
	ctx context.Context,
	req tasks.RawTask,
) (*tasks.Task, error) {
	meta := req.GetHeader()

	taskName := meta.Get("X-CloudTasks-TaskName")
	if taskName == "" {
		return nil, fmt.Errorf("invalid cloud tasks request: missing task name")
	}
	meta.Del("X-CloudTasks-TaskName")

	queueName := meta.Get("X-CloudTasks-QueueName")
	if queueName == "" {
		return nil, fmt.Errorf("invalid cloud tasks request: missing queue name")
	}
	meta.Del("X-CloudTasks-QueueName")

	scheduleTimeValue := meta.Get("X-CloudTasks-TaskETA")
	if scheduleTimeValue == "" {
		return nil, fmt.Errorf("invalid cloud tasks request: missing schedule time")
	}
	meta.Del("X-CloudTasks-TaskETA")

	scheduleTimeDecimal, err := strconv.ParseFloat(scheduleTimeValue, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid cloud tasks request: invalid schedule time: %v", err)
	}
	scheduleTimeSec, scheduleTimeNsec := math.Modf(scheduleTimeDecimal)
	scheduleTime := time.Unix(int64(scheduleTimeSec), int64(scheduleTimeNsec))

	retryCountValue := meta.Get("X-CloudTasks-TaskRetryCount")
	if retryCountValue == "" {
		return nil, fmt.Errorf("invalid cloud tasks request: missing retry count")
	}
	meta.Del("X-CloudTasks-TaskRetryCount")

	retryCount, err := strconv.Atoi(retryCountValue)
	if err != nil {
		return nil, fmt.Errorf("invalid cloud tasks request: invalid retry count: %v", err)
	}

	executionCountValue := meta.Get("X-CloudTasks-TaskExecutionCount")
	if executionCountValue == "" {
		return nil, fmt.Errorf("invalid cloud tasks request: missing execution count")
	}
	meta.Del("X-CloudTasks-TaskExecutionCount")

	executionCount, err := strconv.Atoi(executionCountValue)
	if err != nil {
		return nil, fmt.Errorf("invalid cloud tasks request: invalid execution count: %v", err)
	}

	return &tasks.Task{
		TaskName:       taskName,
		QueueName:      queueName,
		Data:           req.GetData(),
		ScheduleTime:   scheduleTime,
		RetryCount:     retryCount,
		ExecutionCount: executionCount,
		Header:         meta,
	}, nil
}

func (d *CloudTasksDriver) Close() error {
	return d.client.Close()
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

const maxPayloadBytes = int64(65536)

func CloudTasksPushHandler(t tasks.Tasks, h tasks.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxPayloadBytes)

		ctx := r.Context()

		log, err := goji.LoggerFromContext(ctx)
		if err != nil {
			log = zap.NewNop()
		}

		log = log.Named("tasks_push_handler").With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		data, err := tasks.NewPushTaskDataFromRequest(r)
		if err != nil {
			log.Warn("error creating task data", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		task, err := t.Receive(ctx, data)
		if err != nil {
			log.Warn("error recieving task", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log = log.With(
			zap.String("task_name", task.TaskName),
			zap.String("task_queue", task.QueueName),
			zap.Any("task_execution_count", task.ExecutionCount),
		)

		if err := h.HandleTask(ctx, task); err != nil {
			log.Warn("error handling message", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
