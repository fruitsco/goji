package tasks

import (
	"io"
	"net/http"
	"time"
)

type CreateTaskRequest struct {
	// Name is the name of the task
	Name string

	// Data is the payload of the task.
	Data []byte

	// Queue is the name of the queue.
	//  - For cloud tasks, this is the queue name.
	//  - For pubsub, this is the topic name.
	Queue string

	// ScheduleTime is the time the task should be executed.
	// This option is not supported / ignored by the queue driver.
	ScheduleTime *time.Time

	// Url is the URL to send the request to.
	// This option is not supported / ignored by the queue driver.
	Url string

	// Method is the HTTP method to use.
	// This option is not supported / ignored by the queue driver.
	Method string

	// Header are the HTTP headers to send with the request.
	// This option is not supported / ignored by the queue driver.
	Header http.Header
}

type Task struct {
	TaskName       string
	QueueName      string
	ScheduleTime   time.Time
	RetryCount     int
	ExecutionCount int
	Data           []byte
	Header         http.Header
}

type RawTask interface {
	GetData() []byte
	GetHeader() http.Header
}

type PushTaskData struct {
	data   []byte
	header http.Header
}

func (m *PushTaskData) GetData() []byte {
	return m.data
}

func (m *PushTaskData) GetHeader() http.Header {
	return m.header
}

func NewPushTaskData(data []byte, header http.Header) *PushTaskData {
	return &PushTaskData{
		data:   data,
		header: header,
	}
}

func NewPushTaskDataFromRequest(r *http.Request) (*PushTaskData, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return NewPushTaskData(data, r.Header), nil
}
