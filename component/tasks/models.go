package tasks

import "time"

type CreateTaskRequest interface {
	GetName() string
	GetScheduleTime() *time.Time
	isCreateTaskRequest()
}

type CreateHttpTaskRequest struct {
	Name         string
	ScheduleTime *time.Time
	Queue        string
	Url          string
	Method       string
	Headers      map[string]string
	Body         []byte
}

var _ = CreateTaskRequest(&CreateHttpTaskRequest{})

func (r *CreateHttpTaskRequest) GetName() string {
	return r.Name
}

func (r *CreateHttpTaskRequest) GetScheduleTime() *time.Time {
	return r.ScheduleTime
}

func (r *CreateHttpTaskRequest) isCreateTaskRequest() {}

type CreateQueueTaskRequest struct {
	Name         string
	ScheduleTime *time.Time
	Topic        string
	Data         []byte
}

var _ = CreateTaskRequest(&CreateQueueTaskRequest{})

func (r *CreateQueueTaskRequest) GetName() string {
	return r.Name
}

func (r *CreateQueueTaskRequest) GetScheduleTime() *time.Time {
	return r.ScheduleTime
}

func (r *CreateQueueTaskRequest) isCreateTaskRequest() {}

type Task struct {
	TaskName       string
	QueueName      string
	ScheduleTime   time.Time
	RetryCount     int
	ExecutionCount int
	Data           []byte
	Meta           map[string]string
}

type PushRequest struct {
	EndpointName string
	Data         []byte
	Meta         map[string]string
}
