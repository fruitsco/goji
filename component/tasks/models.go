package tasks

import "time"

type CreateTaskRequest interface {
	isCreateTaskRequest()
}

type CreateHttpTaskRequest struct {
	Queue   string
	Url     string
	Method  string
	Headers map[string]string
	Body    []byte
}

var _ = CreateTaskRequest(&CreateHttpTaskRequest{})

func (r *CreateHttpTaskRequest) isCreateTaskRequest() {}

type CreateQueueTaskRequest struct {
	Topic string
	Data  []byte
}

var _ = CreateTaskRequest(&CreateQueueTaskRequest{})

func (r *CreateQueueTaskRequest) isCreateTaskRequest() {}

type Task struct {
	TaskName       string
	ScheduleTime   time.Time
	RetryCount     int
	ExecutionCount int
	Data           []byte
	Meta           map[string]string
}

type PushRequest struct {
	TaskName string
	Data     []byte
	Meta     map[string]string
}
