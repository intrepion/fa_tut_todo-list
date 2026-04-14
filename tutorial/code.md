# Code

### 1. Red: List Tasks From The Store

Create the first code test file:

```bash
mkdir -p workspace/internal/code
touch workspace/internal/code/task_service_test.go
just format
git add --all
git commit --message 'touch workspace/internal/code/task_service_test.go'
```

Put this exact content in `workspace/internal/code/task_service_test.go`:

```go
package code

import (
	"testing"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTaskStore struct {
	mock.Mock
}

func (m *MockTaskStore) ListTasks() ([]contracts.Task, error) {
	args := m.Called()
	return args.Get(0).([]contracts.Task), args.Error(1)
}

func (m *MockTaskStore) CreateTask(taskText string) (contracts.Task, error) {
	args := m.Called(taskText)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskStore) GetTask(taskID int64) (contracts.Task, bool, error) {
	args := m.Called(taskID)
	return args.Get(0).(contracts.Task), args.Bool(1), args.Error(2)
}

func (m *MockTaskStore) DeleteTask(taskID int64) (bool, error) {
	args := m.Called(taskID)
	return args.Bool(0), args.Error(1)
}

func TestTaskServiceListTasksReturnsTasksInOrder(t *testing.T) {
	store := new(MockTaskStore)
	store.On("ListTasks").Return([]contracts.Task{
		{ID: 1, Text: "Learn how to invert binary trees"},
		{ID: 2, Text: "Buy milk"},
	}, nil)

	service := NewTaskService(store)
	result, err := service.ListTasks()

	require.NoError(t, err)
	assert.Equal(t, []contracts.Task{
		{ID: 1, Text: "Learn how to invert binary trees"},
		{ID: 2, Text: "Buy milk"},
	}, result.Tasks)
	store.AssertExpectations(t)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "1. Red: List Tasks From The Store"
```

### 2. Green: Return The Current Task List

Create the first production file:

```bash
mkdir -p workspace/internal/code
touch workspace/internal/code/task_service.go
just format
git add --all
git commit --message 'touch workspace/internal/code/task_service.go'
```

Put this exact content in `workspace/internal/code/task_service.go`:

```go
package code

import "github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"

type DefaultTaskService struct {
	store contracts.TaskStore
}

func NewTaskService(store contracts.TaskStore) contracts.TaskService {
	return DefaultTaskService{store: store}
}

func (s DefaultTaskService) ListTasks() (contracts.TaskListResponse, error) {
	tasks, err := s.store.ListTasks()
	if err != nil {
		return contracts.TaskListResponse{}, err
	}

	return contracts.TaskListResponse{
		Tasks: append([]contracts.Task{}, tasks...),
	}, nil
}

func (s DefaultTaskService) CreateTask(taskText string) (contracts.Task, error) {
	panic("not implemented")
}

func (s DefaultTaskService) GetTask(taskID int64) (contracts.Task, error) {
	panic("not implemented")
}

func (s DefaultTaskService) DeleteTask(taskID int64) error {
	panic("not implemented")
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "2. Green: Return The Current Task List"
```

### 3. Red: Reject Blank Task Creation

Replace `workspace/internal/code/task_service_test.go` with:

```go
package code

import (
	"testing"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTaskStore struct {
	mock.Mock
}

func (m *MockTaskStore) ListTasks() ([]contracts.Task, error) {
	args := m.Called()
	return args.Get(0).([]contracts.Task), args.Error(1)
}

func (m *MockTaskStore) CreateTask(taskText string) (contracts.Task, error) {
	args := m.Called(taskText)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskStore) GetTask(taskID int64) (contracts.Task, bool, error) {
	args := m.Called(taskID)
	return args.Get(0).(contracts.Task), args.Bool(1), args.Error(2)
}

func (m *MockTaskStore) DeleteTask(taskID int64) (bool, error) {
	args := m.Called(taskID)
	return args.Bool(0), args.Error(1)
}

func TestTaskServiceListTasksReturnsTasksInOrder(t *testing.T) {
	store := new(MockTaskStore)
	store.On("ListTasks").Return([]contracts.Task{
		{ID: 1, Text: "Learn how to invert binary trees"},
		{ID: 2, Text: "Buy milk"},
	}, nil)

	service := NewTaskService(store)
	result, err := service.ListTasks()

	require.NoError(t, err)
	assert.Equal(t, []contracts.Task{
		{ID: 1, Text: "Learn how to invert binary trees"},
		{ID: 2, Text: "Buy milk"},
	}, result.Tasks)
	store.AssertExpectations(t)
}

func TestTaskServiceCreateTaskRejectsBlankText(t *testing.T) {
	store := new(MockTaskStore)
	service := NewTaskService(store)

	_, err := service.CreateTask("   ")

	assert.ErrorIs(t, err, contracts.ErrTaskTextBlank)
	store.AssertExpectations(t)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "3. Red: Reject Blank Task Creation"
```

### 4. Green: Create Tasks And Trim Input

Replace `workspace/internal/code/task_service.go` with:

```go
package code

import (
	"strings"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
)

type DefaultTaskService struct {
	store contracts.TaskStore
}

func NewTaskService(store contracts.TaskStore) contracts.TaskService {
	return DefaultTaskService{store: store}
}

func (s DefaultTaskService) ListTasks() (contracts.TaskListResponse, error) {
	tasks, err := s.store.ListTasks()
	if err != nil {
		return contracts.TaskListResponse{}, err
	}

	return contracts.TaskListResponse{
		Tasks: append([]contracts.Task{}, tasks...),
	}, nil
}

func (s DefaultTaskService) CreateTask(taskText string) (contracts.Task, error) {
	trimmed := strings.TrimSpace(taskText)
	if trimmed == "" {
		return contracts.Task{}, contracts.ErrTaskTextBlank
	}

	return s.store.CreateTask(trimmed)
}

func (s DefaultTaskService) GetTask(taskID int64) (contracts.Task, error) {
	panic("not implemented")
}

func (s DefaultTaskService) DeleteTask(taskID int64) error {
	panic("not implemented")
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "4. Green: Create Tasks And Trim Input"
```

### 5. Red: Add Lookup And Delete Not-Found Behavior

Replace `workspace/internal/code/task_service_test.go` with:

```go
package code

import (
	"testing"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTaskStore struct {
	mock.Mock
}

func (m *MockTaskStore) ListTasks() ([]contracts.Task, error) {
	args := m.Called()
	return args.Get(0).([]contracts.Task), args.Error(1)
}

func (m *MockTaskStore) CreateTask(taskText string) (contracts.Task, error) {
	args := m.Called(taskText)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskStore) GetTask(taskID int64) (contracts.Task, bool, error) {
	args := m.Called(taskID)
	return args.Get(0).(contracts.Task), args.Bool(1), args.Error(2)
}

func (m *MockTaskStore) DeleteTask(taskID int64) (bool, error) {
	args := m.Called(taskID)
	return args.Bool(0), args.Error(1)
}

func TestTaskServiceListTasksReturnsTasksInOrder(t *testing.T) {
	store := new(MockTaskStore)
	store.On("ListTasks").Return([]contracts.Task{
		{ID: 1, Text: "Learn how to invert binary trees"},
		{ID: 2, Text: "Buy milk"},
	}, nil)

	service := NewTaskService(store)
	result, err := service.ListTasks()

	require.NoError(t, err)
	assert.Equal(t, []contracts.Task{
		{ID: 1, Text: "Learn how to invert binary trees"},
		{ID: 2, Text: "Buy milk"},
	}, result.Tasks)
	store.AssertExpectations(t)
}

func TestTaskServiceCreateTaskRejectsBlankText(t *testing.T) {
	store := new(MockTaskStore)
	service := NewTaskService(store)

	_, err := service.CreateTask("   ")

	assert.ErrorIs(t, err, contracts.ErrTaskTextBlank)
	store.AssertExpectations(t)
}

func TestTaskServiceGetTaskReturnsNotFoundForMissingIds(t *testing.T) {
	store := new(MockTaskStore)
	store.On("GetTask", int64(9)).Return(contracts.Task{}, false, nil)

	service := NewTaskService(store)
	_, err := service.GetTask(9)

	assert.ErrorIs(t, err, contracts.ErrTaskNotFound)
	store.AssertExpectations(t)
}

func TestTaskServiceDeleteTaskReturnsNotFoundForMissingIds(t *testing.T) {
	store := new(MockTaskStore)
	store.On("DeleteTask", int64(9)).Return(false, nil)

	service := NewTaskService(store)
	err := service.DeleteTask(9)

	assert.ErrorIs(t, err, contracts.ErrTaskNotFound)
	store.AssertExpectations(t)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "5. Red: Add Lookup And Delete Not-Found Behavior"
```

### 6. Green: Return Tasks By Id And Delete Existing Tasks

Replace `workspace/internal/code/task_service.go` with:

```go
package code

import (
	"strings"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
)

type DefaultTaskService struct {
	store contracts.TaskStore
}

func NewTaskService(store contracts.TaskStore) contracts.TaskService {
	return DefaultTaskService{store: store}
}

func (s DefaultTaskService) ListTasks() (contracts.TaskListResponse, error) {
	tasks, err := s.store.ListTasks()
	if err != nil {
		return contracts.TaskListResponse{}, err
	}

	return contracts.TaskListResponse{
		Tasks: append([]contracts.Task{}, tasks...),
	}, nil
}

func (s DefaultTaskService) CreateTask(taskText string) (contracts.Task, error) {
	trimmed := strings.TrimSpace(taskText)
	if trimmed == "" {
		return contracts.Task{}, contracts.ErrTaskTextBlank
	}

	return s.store.CreateTask(trimmed)
}

func (s DefaultTaskService) GetTask(taskID int64) (contracts.Task, error) {
	task, found, err := s.store.GetTask(taskID)
	if err != nil {
		return contracts.Task{}, err
	}
	if !found {
		return contracts.Task{}, contracts.ErrTaskNotFound
	}

	return task, nil
}

func (s DefaultTaskService) DeleteTask(taskID int64) error {
	deleted, err := s.store.DeleteTask(taskID)
	if err != nil {
		return err
	}
	if !deleted {
		return contracts.ErrTaskNotFound
	}

	return nil
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "6. Green: Return Tasks By Id And Delete Existing Tasks"
```
