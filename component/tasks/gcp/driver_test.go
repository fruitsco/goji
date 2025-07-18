package tasksgcp_test

import (
	"context"
	"net/http"
	"testing"

	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	cloudtaskspb_mocks "github.com/fruitsco/goji/test/mocks/googleapis/cloudtaskspb"
	testutil "github.com/fruitsco/goji/test/util"

	"github.com/fruitsco/goji/component/tasks"
	tasksgcp "github.com/fruitsco/goji/component/tasks/gcp"
)

func createMockServer(t *testing.T) (*grpc.ClientConn, *cloudtaskspb_mocks.MockCloudTasksServer) {
	// create the cloudtasks mock server
	mockServer := cloudtaskspb_mocks.NewMockCloudTasksServer(t)

	_, c := testutil.NewGrpcServer(t, func(s *grpc.Server) {
		cloudtaskspb.RegisterCloudTasksServer(s, mockServer)
	})

	return c, mockServer
}

func TestCloudTasksDriver_Submit(t *testing.T) {
	client, mockServer := createMockServer(t)

	mockServer.EXPECT().CreateTask(mock.Anything, mock.Anything).RunAndReturn(func(_ context.Context, ctr *cloudtaskspb.CreateTaskRequest) (*cloudtaskspb.Task, error) {
		assert.Equal(t, "projects/test-project/locations/test-region/queues/test-queue", ctr.Parent)
		assert.Equal(t, "projects/test-project/locations/test-region/queues/test-queue/tasks/test-task", ctr.Task.Name)
		assert.Equal(t, "http://test.local", ctr.Task.GetHttpRequest().GetUrl())
		assert.Equal(t, cloudtaskspb.HttpMethod_POST, ctr.Task.GetHttpRequest().GetHttpMethod())
		assert.Equal(t, []byte("test"), ctr.Task.GetHttpRequest().GetBody())
		return &cloudtaskspb.Task{}, nil
	})

	driver, err := tasksgcp.NewCloudTasksDriver(tasksgcp.CloudTasksDriverParams{
		Context: context.Background(),
		Config: &tasks.CloudTasksConfig{
			ProjectID:  "test-project",
			Region:     "test-region",
			DefaultUrl: "http://test.local",
		},
		NoAuth:   true,
		GRPCConn: client,
		Log:      zap.NewNop(),
	})
	require.NoError(t, err)
	defer driver.Close()

	err = driver.Submit(context.Background(), &tasks.CreateTaskRequest{
		Name:  "test-task",
		Queue: "test-queue",
		Data:  []byte("test"),
	})
	require.NoError(t, err)
}

func TestCloudTasksDriver_ReceivePush(t *testing.T) {
	driver, err := tasksgcp.NewCloudTasksDriver(tasksgcp.CloudTasksDriverParams{
		Context: context.Background(),
		Config: &tasks.CloudTasksConfig{
			ProjectID:  "test-project",
			Region:     "test-region",
			DefaultUrl: "http://test.local",
		},
		NoAuth: true,
		Log:    zap.NewNop(),
	})
	require.NoError(t, err)
	defer driver.Close()

	req := tasks.NewPushTaskData(
		[]byte("test"),
		http.Header{
			"X-Cloudtasks-Taskname":           []string{"test-task"},
			"X-Cloudtasks-Queuename":          []string{"test-queue"},
			"X-Cloudtasks-Tasketa":            []string{"1723123865.123456789"},
			"X-Cloudtasks-Taskretrycount":     []string{"0"},
			"X-Cloudtasks-Taskexecutioncount": []string{"1"},
		},
	)

	task, err := driver.Receive(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "test-task", task.TaskName)
	assert.Equal(t, "test-queue", task.QueueName)
	assert.Equal(t, int64(1723123865), task.ScheduleTime.Unix())
	assert.Equal(t, 0, task.RetryCount)
	assert.Equal(t, 1, task.ExecutionCount)
}

func TestCloudTasksDriver_ReceivePush_FailsForInvalidHeaders(t *testing.T) {
	driver, err := tasksgcp.NewCloudTasksDriver(tasksgcp.CloudTasksDriverParams{
		Context: context.Background(),
		Config: &tasks.CloudTasksConfig{
			ProjectID:  "test-project",
			Region:     "test-region",
			DefaultUrl: "http://test.local",
		},
		NoAuth: true,
		Log:    zap.NewNop(),
	})
	require.NoError(t, err)
	defer driver.Close()

	validHeaders := http.Header{
		"X-Cloudtasks-Taskname":           []string{"test-task"},
		"X-Cloudtasks-Queuename":          []string{"test-queue"},
		"X-Cloudtasks-Tasketa":            []string{"1723123865.123456789"},
		"X-Cloudtasks-Taskretrycount":     []string{"0"},
		"X-Cloudtasks-Taskexecutioncount": []string{"1"},
	}

	for header := range validHeaders {
		t.Run("missing "+header, func(t *testing.T) {
			invalidHeaders := validHeaders.Clone()
			invalidHeaders.Del(header)

			req := tasks.NewPushTaskData(nil, invalidHeaders)

			_, err := driver.Receive(context.Background(), req)
			require.Error(t, err)
		})
	}
}
