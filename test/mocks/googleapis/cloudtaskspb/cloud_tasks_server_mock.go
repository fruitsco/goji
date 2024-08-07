// Code generated by mockery v2.43.2. DO NOT EDIT.

package cloudtaskspb_mocks

import (
	context "context"

	cloudtaskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	iampb "cloud.google.com/go/iam/apiv1/iampb"

	mock "github.com/stretchr/testify/mock"
)

// MockCloudTasksServer is an autogenerated mock type for the CloudTasksServer type
type MockCloudTasksServer struct {
	mock.Mock
}

type MockCloudTasksServer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCloudTasksServer) EXPECT() *MockCloudTasksServer_Expecter {
	return &MockCloudTasksServer_Expecter{mock: &_m.Mock}
}

// CreateQueue provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) CreateQueue(_a0 context.Context, _a1 *cloudtaskspb.CreateQueueRequest) (*cloudtaskspb.Queue, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateQueue")
	}

	var r0 *cloudtaskspb.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.CreateQueueRequest) (*cloudtaskspb.Queue, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.CreateQueueRequest) *cloudtaskspb.Queue); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Queue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.CreateQueueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_CreateQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateQueue'
type MockCloudTasksServer_CreateQueue_Call struct {
	*mock.Call
}

// CreateQueue is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.CreateQueueRequest
func (_e *MockCloudTasksServer_Expecter) CreateQueue(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_CreateQueue_Call {
	return &MockCloudTasksServer_CreateQueue_Call{Call: _e.mock.On("CreateQueue", _a0, _a1)}
}

func (_c *MockCloudTasksServer_CreateQueue_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.CreateQueueRequest)) *MockCloudTasksServer_CreateQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.CreateQueueRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_CreateQueue_Call) Return(_a0 *cloudtaskspb.Queue, _a1 error) *MockCloudTasksServer_CreateQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_CreateQueue_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.CreateQueueRequest) (*cloudtaskspb.Queue, error)) *MockCloudTasksServer_CreateQueue_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTask provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) CreateTask(_a0 context.Context, _a1 *cloudtaskspb.CreateTaskRequest) (*cloudtaskspb.Task, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateTask")
	}

	var r0 *cloudtaskspb.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.CreateTaskRequest) (*cloudtaskspb.Task, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.CreateTaskRequest) *cloudtaskspb.Task); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.CreateTaskRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_CreateTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTask'
type MockCloudTasksServer_CreateTask_Call struct {
	*mock.Call
}

// CreateTask is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.CreateTaskRequest
func (_e *MockCloudTasksServer_Expecter) CreateTask(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_CreateTask_Call {
	return &MockCloudTasksServer_CreateTask_Call{Call: _e.mock.On("CreateTask", _a0, _a1)}
}

func (_c *MockCloudTasksServer_CreateTask_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.CreateTaskRequest)) *MockCloudTasksServer_CreateTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.CreateTaskRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_CreateTask_Call) Return(_a0 *cloudtaskspb.Task, _a1 error) *MockCloudTasksServer_CreateTask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_CreateTask_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.CreateTaskRequest) (*cloudtaskspb.Task, error)) *MockCloudTasksServer_CreateTask_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteQueue provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) DeleteQueue(_a0 context.Context, _a1 *cloudtaskspb.DeleteQueueRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteQueue")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.DeleteQueueRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.DeleteQueueRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.DeleteQueueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_DeleteQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteQueue'
type MockCloudTasksServer_DeleteQueue_Call struct {
	*mock.Call
}

// DeleteQueue is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.DeleteQueueRequest
func (_e *MockCloudTasksServer_Expecter) DeleteQueue(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_DeleteQueue_Call {
	return &MockCloudTasksServer_DeleteQueue_Call{Call: _e.mock.On("DeleteQueue", _a0, _a1)}
}

func (_c *MockCloudTasksServer_DeleteQueue_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.DeleteQueueRequest)) *MockCloudTasksServer_DeleteQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.DeleteQueueRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_DeleteQueue_Call) Return(_a0 *emptypb.Empty, _a1 error) *MockCloudTasksServer_DeleteQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_DeleteQueue_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.DeleteQueueRequest) (*emptypb.Empty, error)) *MockCloudTasksServer_DeleteQueue_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTask provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) DeleteTask(_a0 context.Context, _a1 *cloudtaskspb.DeleteTaskRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTask")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.DeleteTaskRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.DeleteTaskRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.DeleteTaskRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_DeleteTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTask'
type MockCloudTasksServer_DeleteTask_Call struct {
	*mock.Call
}

// DeleteTask is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.DeleteTaskRequest
func (_e *MockCloudTasksServer_Expecter) DeleteTask(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_DeleteTask_Call {
	return &MockCloudTasksServer_DeleteTask_Call{Call: _e.mock.On("DeleteTask", _a0, _a1)}
}

func (_c *MockCloudTasksServer_DeleteTask_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.DeleteTaskRequest)) *MockCloudTasksServer_DeleteTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.DeleteTaskRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_DeleteTask_Call) Return(_a0 *emptypb.Empty, _a1 error) *MockCloudTasksServer_DeleteTask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_DeleteTask_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.DeleteTaskRequest) (*emptypb.Empty, error)) *MockCloudTasksServer_DeleteTask_Call {
	_c.Call.Return(run)
	return _c
}

// GetIamPolicy provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) GetIamPolicy(_a0 context.Context, _a1 *iampb.GetIamPolicyRequest) (*iampb.Policy, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetIamPolicy")
	}

	var r0 *iampb.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iampb.GetIamPolicyRequest) (*iampb.Policy, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iampb.GetIamPolicyRequest) *iampb.Policy); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iampb.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iampb.GetIamPolicyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_GetIamPolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIamPolicy'
type MockCloudTasksServer_GetIamPolicy_Call struct {
	*mock.Call
}

// GetIamPolicy is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *iampb.GetIamPolicyRequest
func (_e *MockCloudTasksServer_Expecter) GetIamPolicy(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_GetIamPolicy_Call {
	return &MockCloudTasksServer_GetIamPolicy_Call{Call: _e.mock.On("GetIamPolicy", _a0, _a1)}
}

func (_c *MockCloudTasksServer_GetIamPolicy_Call) Run(run func(_a0 context.Context, _a1 *iampb.GetIamPolicyRequest)) *MockCloudTasksServer_GetIamPolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*iampb.GetIamPolicyRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_GetIamPolicy_Call) Return(_a0 *iampb.Policy, _a1 error) *MockCloudTasksServer_GetIamPolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_GetIamPolicy_Call) RunAndReturn(run func(context.Context, *iampb.GetIamPolicyRequest) (*iampb.Policy, error)) *MockCloudTasksServer_GetIamPolicy_Call {
	_c.Call.Return(run)
	return _c
}

// GetQueue provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) GetQueue(_a0 context.Context, _a1 *cloudtaskspb.GetQueueRequest) (*cloudtaskspb.Queue, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetQueue")
	}

	var r0 *cloudtaskspb.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.GetQueueRequest) (*cloudtaskspb.Queue, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.GetQueueRequest) *cloudtaskspb.Queue); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Queue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.GetQueueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_GetQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQueue'
type MockCloudTasksServer_GetQueue_Call struct {
	*mock.Call
}

// GetQueue is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.GetQueueRequest
func (_e *MockCloudTasksServer_Expecter) GetQueue(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_GetQueue_Call {
	return &MockCloudTasksServer_GetQueue_Call{Call: _e.mock.On("GetQueue", _a0, _a1)}
}

func (_c *MockCloudTasksServer_GetQueue_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.GetQueueRequest)) *MockCloudTasksServer_GetQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.GetQueueRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_GetQueue_Call) Return(_a0 *cloudtaskspb.Queue, _a1 error) *MockCloudTasksServer_GetQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_GetQueue_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.GetQueueRequest) (*cloudtaskspb.Queue, error)) *MockCloudTasksServer_GetQueue_Call {
	_c.Call.Return(run)
	return _c
}

// GetTask provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) GetTask(_a0 context.Context, _a1 *cloudtaskspb.GetTaskRequest) (*cloudtaskspb.Task, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetTask")
	}

	var r0 *cloudtaskspb.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.GetTaskRequest) (*cloudtaskspb.Task, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.GetTaskRequest) *cloudtaskspb.Task); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.GetTaskRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_GetTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTask'
type MockCloudTasksServer_GetTask_Call struct {
	*mock.Call
}

// GetTask is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.GetTaskRequest
func (_e *MockCloudTasksServer_Expecter) GetTask(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_GetTask_Call {
	return &MockCloudTasksServer_GetTask_Call{Call: _e.mock.On("GetTask", _a0, _a1)}
}

func (_c *MockCloudTasksServer_GetTask_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.GetTaskRequest)) *MockCloudTasksServer_GetTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.GetTaskRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_GetTask_Call) Return(_a0 *cloudtaskspb.Task, _a1 error) *MockCloudTasksServer_GetTask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_GetTask_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.GetTaskRequest) (*cloudtaskspb.Task, error)) *MockCloudTasksServer_GetTask_Call {
	_c.Call.Return(run)
	return _c
}

// ListQueues provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) ListQueues(_a0 context.Context, _a1 *cloudtaskspb.ListQueuesRequest) (*cloudtaskspb.ListQueuesResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListQueues")
	}

	var r0 *cloudtaskspb.ListQueuesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.ListQueuesRequest) (*cloudtaskspb.ListQueuesResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.ListQueuesRequest) *cloudtaskspb.ListQueuesResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.ListQueuesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.ListQueuesRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_ListQueues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListQueues'
type MockCloudTasksServer_ListQueues_Call struct {
	*mock.Call
}

// ListQueues is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.ListQueuesRequest
func (_e *MockCloudTasksServer_Expecter) ListQueues(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_ListQueues_Call {
	return &MockCloudTasksServer_ListQueues_Call{Call: _e.mock.On("ListQueues", _a0, _a1)}
}

func (_c *MockCloudTasksServer_ListQueues_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.ListQueuesRequest)) *MockCloudTasksServer_ListQueues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.ListQueuesRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_ListQueues_Call) Return(_a0 *cloudtaskspb.ListQueuesResponse, _a1 error) *MockCloudTasksServer_ListQueues_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_ListQueues_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.ListQueuesRequest) (*cloudtaskspb.ListQueuesResponse, error)) *MockCloudTasksServer_ListQueues_Call {
	_c.Call.Return(run)
	return _c
}

// ListTasks provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) ListTasks(_a0 context.Context, _a1 *cloudtaskspb.ListTasksRequest) (*cloudtaskspb.ListTasksResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListTasks")
	}

	var r0 *cloudtaskspb.ListTasksResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.ListTasksRequest) (*cloudtaskspb.ListTasksResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.ListTasksRequest) *cloudtaskspb.ListTasksResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.ListTasksResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.ListTasksRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_ListTasks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListTasks'
type MockCloudTasksServer_ListTasks_Call struct {
	*mock.Call
}

// ListTasks is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.ListTasksRequest
func (_e *MockCloudTasksServer_Expecter) ListTasks(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_ListTasks_Call {
	return &MockCloudTasksServer_ListTasks_Call{Call: _e.mock.On("ListTasks", _a0, _a1)}
}

func (_c *MockCloudTasksServer_ListTasks_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.ListTasksRequest)) *MockCloudTasksServer_ListTasks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.ListTasksRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_ListTasks_Call) Return(_a0 *cloudtaskspb.ListTasksResponse, _a1 error) *MockCloudTasksServer_ListTasks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_ListTasks_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.ListTasksRequest) (*cloudtaskspb.ListTasksResponse, error)) *MockCloudTasksServer_ListTasks_Call {
	_c.Call.Return(run)
	return _c
}

// PauseQueue provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) PauseQueue(_a0 context.Context, _a1 *cloudtaskspb.PauseQueueRequest) (*cloudtaskspb.Queue, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PauseQueue")
	}

	var r0 *cloudtaskspb.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.PauseQueueRequest) (*cloudtaskspb.Queue, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.PauseQueueRequest) *cloudtaskspb.Queue); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Queue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.PauseQueueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_PauseQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PauseQueue'
type MockCloudTasksServer_PauseQueue_Call struct {
	*mock.Call
}

// PauseQueue is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.PauseQueueRequest
func (_e *MockCloudTasksServer_Expecter) PauseQueue(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_PauseQueue_Call {
	return &MockCloudTasksServer_PauseQueue_Call{Call: _e.mock.On("PauseQueue", _a0, _a1)}
}

func (_c *MockCloudTasksServer_PauseQueue_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.PauseQueueRequest)) *MockCloudTasksServer_PauseQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.PauseQueueRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_PauseQueue_Call) Return(_a0 *cloudtaskspb.Queue, _a1 error) *MockCloudTasksServer_PauseQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_PauseQueue_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.PauseQueueRequest) (*cloudtaskspb.Queue, error)) *MockCloudTasksServer_PauseQueue_Call {
	_c.Call.Return(run)
	return _c
}

// PurgeQueue provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) PurgeQueue(_a0 context.Context, _a1 *cloudtaskspb.PurgeQueueRequest) (*cloudtaskspb.Queue, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PurgeQueue")
	}

	var r0 *cloudtaskspb.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.PurgeQueueRequest) (*cloudtaskspb.Queue, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.PurgeQueueRequest) *cloudtaskspb.Queue); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Queue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.PurgeQueueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_PurgeQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PurgeQueue'
type MockCloudTasksServer_PurgeQueue_Call struct {
	*mock.Call
}

// PurgeQueue is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.PurgeQueueRequest
func (_e *MockCloudTasksServer_Expecter) PurgeQueue(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_PurgeQueue_Call {
	return &MockCloudTasksServer_PurgeQueue_Call{Call: _e.mock.On("PurgeQueue", _a0, _a1)}
}

func (_c *MockCloudTasksServer_PurgeQueue_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.PurgeQueueRequest)) *MockCloudTasksServer_PurgeQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.PurgeQueueRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_PurgeQueue_Call) Return(_a0 *cloudtaskspb.Queue, _a1 error) *MockCloudTasksServer_PurgeQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_PurgeQueue_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.PurgeQueueRequest) (*cloudtaskspb.Queue, error)) *MockCloudTasksServer_PurgeQueue_Call {
	_c.Call.Return(run)
	return _c
}

// ResumeQueue provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) ResumeQueue(_a0 context.Context, _a1 *cloudtaskspb.ResumeQueueRequest) (*cloudtaskspb.Queue, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ResumeQueue")
	}

	var r0 *cloudtaskspb.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.ResumeQueueRequest) (*cloudtaskspb.Queue, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.ResumeQueueRequest) *cloudtaskspb.Queue); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Queue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.ResumeQueueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_ResumeQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResumeQueue'
type MockCloudTasksServer_ResumeQueue_Call struct {
	*mock.Call
}

// ResumeQueue is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.ResumeQueueRequest
func (_e *MockCloudTasksServer_Expecter) ResumeQueue(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_ResumeQueue_Call {
	return &MockCloudTasksServer_ResumeQueue_Call{Call: _e.mock.On("ResumeQueue", _a0, _a1)}
}

func (_c *MockCloudTasksServer_ResumeQueue_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.ResumeQueueRequest)) *MockCloudTasksServer_ResumeQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.ResumeQueueRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_ResumeQueue_Call) Return(_a0 *cloudtaskspb.Queue, _a1 error) *MockCloudTasksServer_ResumeQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_ResumeQueue_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.ResumeQueueRequest) (*cloudtaskspb.Queue, error)) *MockCloudTasksServer_ResumeQueue_Call {
	_c.Call.Return(run)
	return _c
}

// RunTask provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) RunTask(_a0 context.Context, _a1 *cloudtaskspb.RunTaskRequest) (*cloudtaskspb.Task, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for RunTask")
	}

	var r0 *cloudtaskspb.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.RunTaskRequest) (*cloudtaskspb.Task, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.RunTaskRequest) *cloudtaskspb.Task); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.RunTaskRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_RunTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RunTask'
type MockCloudTasksServer_RunTask_Call struct {
	*mock.Call
}

// RunTask is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.RunTaskRequest
func (_e *MockCloudTasksServer_Expecter) RunTask(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_RunTask_Call {
	return &MockCloudTasksServer_RunTask_Call{Call: _e.mock.On("RunTask", _a0, _a1)}
}

func (_c *MockCloudTasksServer_RunTask_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.RunTaskRequest)) *MockCloudTasksServer_RunTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.RunTaskRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_RunTask_Call) Return(_a0 *cloudtaskspb.Task, _a1 error) *MockCloudTasksServer_RunTask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_RunTask_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.RunTaskRequest) (*cloudtaskspb.Task, error)) *MockCloudTasksServer_RunTask_Call {
	_c.Call.Return(run)
	return _c
}

// SetIamPolicy provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) SetIamPolicy(_a0 context.Context, _a1 *iampb.SetIamPolicyRequest) (*iampb.Policy, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetIamPolicy")
	}

	var r0 *iampb.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iampb.SetIamPolicyRequest) (*iampb.Policy, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iampb.SetIamPolicyRequest) *iampb.Policy); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iampb.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iampb.SetIamPolicyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_SetIamPolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetIamPolicy'
type MockCloudTasksServer_SetIamPolicy_Call struct {
	*mock.Call
}

// SetIamPolicy is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *iampb.SetIamPolicyRequest
func (_e *MockCloudTasksServer_Expecter) SetIamPolicy(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_SetIamPolicy_Call {
	return &MockCloudTasksServer_SetIamPolicy_Call{Call: _e.mock.On("SetIamPolicy", _a0, _a1)}
}

func (_c *MockCloudTasksServer_SetIamPolicy_Call) Run(run func(_a0 context.Context, _a1 *iampb.SetIamPolicyRequest)) *MockCloudTasksServer_SetIamPolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*iampb.SetIamPolicyRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_SetIamPolicy_Call) Return(_a0 *iampb.Policy, _a1 error) *MockCloudTasksServer_SetIamPolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_SetIamPolicy_Call) RunAndReturn(run func(context.Context, *iampb.SetIamPolicyRequest) (*iampb.Policy, error)) *MockCloudTasksServer_SetIamPolicy_Call {
	_c.Call.Return(run)
	return _c
}

// TestIamPermissions provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) TestIamPermissions(_a0 context.Context, _a1 *iampb.TestIamPermissionsRequest) (*iampb.TestIamPermissionsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for TestIamPermissions")
	}

	var r0 *iampb.TestIamPermissionsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iampb.TestIamPermissionsRequest) (*iampb.TestIamPermissionsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iampb.TestIamPermissionsRequest) *iampb.TestIamPermissionsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iampb.TestIamPermissionsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iampb.TestIamPermissionsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_TestIamPermissions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TestIamPermissions'
type MockCloudTasksServer_TestIamPermissions_Call struct {
	*mock.Call
}

// TestIamPermissions is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *iampb.TestIamPermissionsRequest
func (_e *MockCloudTasksServer_Expecter) TestIamPermissions(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_TestIamPermissions_Call {
	return &MockCloudTasksServer_TestIamPermissions_Call{Call: _e.mock.On("TestIamPermissions", _a0, _a1)}
}

func (_c *MockCloudTasksServer_TestIamPermissions_Call) Run(run func(_a0 context.Context, _a1 *iampb.TestIamPermissionsRequest)) *MockCloudTasksServer_TestIamPermissions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*iampb.TestIamPermissionsRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_TestIamPermissions_Call) Return(_a0 *iampb.TestIamPermissionsResponse, _a1 error) *MockCloudTasksServer_TestIamPermissions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_TestIamPermissions_Call) RunAndReturn(run func(context.Context, *iampb.TestIamPermissionsRequest) (*iampb.TestIamPermissionsResponse, error)) *MockCloudTasksServer_TestIamPermissions_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateQueue provides a mock function with given fields: _a0, _a1
func (_m *MockCloudTasksServer) UpdateQueue(_a0 context.Context, _a1 *cloudtaskspb.UpdateQueueRequest) (*cloudtaskspb.Queue, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateQueue")
	}

	var r0 *cloudtaskspb.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.UpdateQueueRequest) (*cloudtaskspb.Queue, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtaskspb.UpdateQueueRequest) *cloudtaskspb.Queue); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtaskspb.Queue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cloudtaskspb.UpdateQueueRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloudTasksServer_UpdateQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateQueue'
type MockCloudTasksServer_UpdateQueue_Call struct {
	*mock.Call
}

// UpdateQueue is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cloudtaskspb.UpdateQueueRequest
func (_e *MockCloudTasksServer_Expecter) UpdateQueue(_a0 interface{}, _a1 interface{}) *MockCloudTasksServer_UpdateQueue_Call {
	return &MockCloudTasksServer_UpdateQueue_Call{Call: _e.mock.On("UpdateQueue", _a0, _a1)}
}

func (_c *MockCloudTasksServer_UpdateQueue_Call) Run(run func(_a0 context.Context, _a1 *cloudtaskspb.UpdateQueueRequest)) *MockCloudTasksServer_UpdateQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cloudtaskspb.UpdateQueueRequest))
	})
	return _c
}

func (_c *MockCloudTasksServer_UpdateQueue_Call) Return(_a0 *cloudtaskspb.Queue, _a1 error) *MockCloudTasksServer_UpdateQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloudTasksServer_UpdateQueue_Call) RunAndReturn(run func(context.Context, *cloudtaskspb.UpdateQueueRequest) (*cloudtaskspb.Queue, error)) *MockCloudTasksServer_UpdateQueue_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCloudTasksServer creates a new instance of MockCloudTasksServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCloudTasksServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCloudTasksServer {
	mock := &MockCloudTasksServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}