package tasks

import "time"

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

	// Headers are the HTTP headers to send with the request.
	// This option is not supported / ignored by the queue driver.
	Headers map[string]string
}

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
	Data []byte
	Meta map[string]string
}
