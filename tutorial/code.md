# Code

### 1. Red: Parse The Canonical Stored Task Data

Create the first code test file:

```bash
touch workspace/internal/code/task_list_service_test.go
just format
git add --all
git commit --message 'touch workspace/internal/code/task_list_service_test.go'
```

Put this exact content in `workspace/internal/code/task_list_service_test.go`:

```go
package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTaskStoragePreservesTheCanonicalTasksInOrder(t *testing.T) {
	storageText := "[\n  \"Learn how to invert binary trees\",\n  \"Buy milk\",\n  \"Clean kitchen\"\n]"

	result := parseTaskStorage(storageText)

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, result)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "1. Red: Parse The Canonical Stored Task Data"
```

### 2. Green: Parse The Canonical Stored Task Data

Create the first production file:

```bash
touch workspace/internal/code/task_list_service.go
just format
git add --all
git commit --message 'touch workspace/internal/code/task_list_service.go'
```

Put this exact content in `workspace/internal/code/task_list_service.go`:

```go
package code

import "encoding/json"

func parseTaskStorage(storageText string) []string {
	var tasks []string
	_ = json.Unmarshal([]byte(storageText), &tasks)
	return tasks
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "2. Green: Parse The Canonical Stored Task Data"
```

### 3. Red: Append And Remove Tasks

Replace `workspace/internal/code/task_list_service_test.go` with:

```go
package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTaskStoragePreservesTheCanonicalTasksInOrder(t *testing.T) {
	storageText := "[\n  \"Learn how to invert binary trees\",\n  \"Buy milk\",\n  \"Clean kitchen\"\n]"

	result := parseTaskStorage(storageText)

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, result)
}

func TestAppendTaskAppendsTheTaskToTheEndOfANewList(t *testing.T) {
	result := appendTask([]string{
		"Learn how to invert binary trees",
		"Buy milk",
	}, "Clean kitchen")

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, result)
}

func TestRemoveTaskByExactTextRemovesTheFirstExactMatch(t *testing.T) {
	result := removeTaskByExactText([]string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, "Buy milk")

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Clean kitchen",
	}, result)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "3. Red: Append And Remove Tasks"
```

### 4. Green: Append And Remove Tasks

Replace `workspace/internal/code/task_list_service.go` with:

```go
package code

import "encoding/json"

func parseTaskStorage(storageText string) []string {
	var tasks []string
	_ = json.Unmarshal([]byte(storageText), &tasks)
	return tasks
}

func appendTask(taskList []string, taskText string) []string {
	next := append([]string{}, taskList...)
	return append(next, taskText)
}

func removeTaskByExactText(taskList []string, completedTaskText string) []string {
	next := append([]string{}, taskList...)
	for index, task := range next {
		if task == completedTaskText {
			return append(next[:index], next[index+1:]...)
		}
	}
	return next
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "4. Green: Append And Remove Tasks"
```

### 5. Red: Format And Serialize The Task List

Replace `workspace/internal/code/task_list_service_test.go` with:

```go
package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTaskStoragePreservesTheCanonicalTasksInOrder(t *testing.T) {
	storageText := "[\n  \"Learn how to invert binary trees\",\n  \"Buy milk\",\n  \"Clean kitchen\"\n]"

	result := parseTaskStorage(storageText)

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, result)
}

func TestAppendTaskAppendsTheTaskToTheEndOfANewList(t *testing.T) {
	result := appendTask([]string{
		"Learn how to invert binary trees",
		"Buy milk",
	}, "Clean kitchen")

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, result)
}

func TestRemoveTaskByExactTextRemovesTheFirstExactMatch(t *testing.T) {
	result := removeTaskByExactText([]string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, "Buy milk")

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Clean kitchen",
	}, result)
}

func TestFormatTaskListReturnsOneLinePerTaskInOrder(t *testing.T) {
	result := formatTaskList([]string{
		"Learn how to invert binary trees",
		"Buy milk",
	})

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
	}, result)
}

func TestSerializeTaskStorageReturnsAJsonArrayInOrder(t *testing.T) {
	result := serializeTaskStorage([]string{
		"Learn how to invert binary trees",
		"Buy milk",
	})

	assert.JSONEq(t, "[\"Learn how to invert binary trees\",\"Buy milk\"]", result)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "5. Red: Format And Serialize The Task List"
```

### 6. Green: Format And Serialize The Task List

Replace `workspace/internal/code/task_list_service.go` with:

```go
package code

import "encoding/json"

func parseTaskStorage(storageText string) []string {
	var tasks []string
	_ = json.Unmarshal([]byte(storageText), &tasks)
	return tasks
}

func appendTask(taskList []string, taskText string) []string {
	next := append([]string{}, taskList...)
	return append(next, taskText)
}

func removeTaskByExactText(taskList []string, completedTaskText string) []string {
	next := append([]string{}, taskList...)
	for index, task := range next {
		if task == completedTaskText {
			return append(next[:index], next[index+1:]...)
		}
	}
	return next
}

func formatTaskList(taskList []string) []string {
	return append([]string{}, taskList...)
}

func serializeTaskStorage(taskList []string) string {
	storageBytes, _ := json.Marshal(taskList)
	return string(storageBytes)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "6. Green: Format And Serialize The Task List"
```

### 7. Red: Add The Task Service Tests

Replace `workspace/internal/code/task_list_service_test.go` with:

```go
package code

import (
	"testing"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskStore struct {
	mock.Mock
}

func (m *MockTaskStore) LoadTaskStorage() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockTaskStore) SaveTaskStorage(storageText string) error {
	args := m.Called(storageText)
	return args.Error(0)
}

func TestParseTaskStoragePreservesTheCanonicalTasksInOrder(t *testing.T) {
	storageText := "[\n  \"Learn how to invert binary trees\",\n  \"Buy milk\",\n  \"Clean kitchen\"\n]"

	result := parseTaskStorage(storageText)

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, result)
}

func TestAppendTaskAppendsTheTaskToTheEndOfANewList(t *testing.T) {
	result := appendTask([]string{
		"Learn how to invert binary trees",
		"Buy milk",
	}, "Clean kitchen")

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, result)
}

func TestRemoveTaskByExactTextRemovesTheFirstExactMatch(t *testing.T) {
	result := removeTaskByExactText([]string{
		"Learn how to invert binary trees",
		"Buy milk",
		"Clean kitchen",
	}, "Buy milk")

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Clean kitchen",
	}, result)
}

func TestFormatTaskListReturnsOneLinePerTaskInOrder(t *testing.T) {
	result := formatTaskList([]string{
		"Learn how to invert binary trees",
		"Buy milk",
	})

	assert.Equal(t, []string{
		"Learn how to invert binary trees",
		"Buy milk",
	}, result)
}

func TestSerializeTaskStorageReturnsAJsonArrayInOrder(t *testing.T) {
	result := serializeTaskStorage([]string{
		"Learn how to invert binary trees",
		"Buy milk",
	})

	assert.JSONEq(t, "[\"Learn how to invert binary trees\",\"Buy milk\"]", result)
}

func TestTaskListServiceListTasksLoadsAndFormatsTasks(t *testing.T) {
	store := new(MockTaskStore)
	store.On("LoadTaskStorage").Return("[\"Learn how to invert binary trees\",\"Buy milk\"]", nil)

	service := NewTaskListService(store)
	result, err := service.ListTasks()

	assert.NoError(t, err)
	assert.Equal(t, contracts.TaskListResult{
		Tasks: []string{"Learn how to invert binary trees", "Buy milk"},
		Lines: []string{"Learn how to invert binary trees", "Buy milk"},
	}, result)
	store.AssertExpectations(t)
}

func TestTaskListServiceAddTaskPersistsUpdatedTasks(t *testing.T) {
	store := new(MockTaskStore)
	store.On("LoadTaskStorage").Return("[\"Learn how to invert binary trees\"]", nil)
	store.On("SaveTaskStorage", "[\"Learn how to invert binary trees\",\"Buy milk\"]").Return(nil)

	service := NewTaskListService(store)
	result, err := service.AddTask("Buy milk")

	assert.NoError(t, err)
	assert.Equal(t, contracts.TaskListResult{
		Tasks: []string{"Learn how to invert binary trees", "Buy milk"},
		Lines: []string{"Learn how to invert binary trees", "Buy milk"},
	}, result)
	store.AssertExpectations(t)
}

func TestTaskListServiceRemoveTaskPersistsUpdatedTasks(t *testing.T) {
	store := new(MockTaskStore)
	store.On("LoadTaskStorage").Return("[\"Learn how to invert binary trees\",\"Buy milk\"]", nil)
	store.On("SaveTaskStorage", "[\"Learn how to invert binary trees\"]").Return(nil)

	service := NewTaskListService(store)
	result, err := service.RemoveTask("Buy milk")

	assert.NoError(t, err)
	assert.Equal(t, contracts.TaskListResult{
		Tasks: []string{"Learn how to invert binary trees"},
		Lines: []string{"Learn how to invert binary trees"},
	}, result)
	store.AssertExpectations(t)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "7. Red: Add The Task Service Tests"
```

### 8. Green: Add The Task Service

Replace `workspace/internal/code/task_list_service.go` with:

```go
package code

import (
	"encoding/json"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
)

type DefaultTaskListService struct {
	store contracts.TaskStore
}

func NewTaskListService(store contracts.TaskStore) contracts.TaskListService {
	return DefaultTaskListService{store: store}
}

func parseTaskStorage(storageText string) []string {
	var tasks []string
	_ = json.Unmarshal([]byte(storageText), &tasks)
	return tasks
}

func appendTask(taskList []string, taskText string) []string {
	next := append([]string{}, taskList...)
	return append(next, taskText)
}

func removeTaskByExactText(taskList []string, completedTaskText string) []string {
	next := append([]string{}, taskList...)
	for index, task := range next {
		if task == completedTaskText {
			return append(next[:index], next[index+1:]...)
		}
	}
	return next
}

func formatTaskList(taskList []string) []string {
	return append([]string{}, taskList...)
}

func serializeTaskStorage(taskList []string) string {
	storageBytes, _ := json.Marshal(taskList)
	return string(storageBytes)
}

func (s DefaultTaskListService) ListTasks() (contracts.TaskListResult, error) {
	storageText, err := s.store.LoadTaskStorage()
	if err != nil {
		return contracts.TaskListResult{}, err
	}

	tasks := parseTaskStorage(storageText)
	return buildTaskListResult(tasks), nil
}

func (s DefaultTaskListService) AddTask(taskText string) (contracts.TaskListResult, error) {
	storageText, err := s.store.LoadTaskStorage()
	if err != nil {
		return contracts.TaskListResult{}, err
	}

	nextTasks := appendTask(parseTaskStorage(storageText), taskText)
	if err := s.store.SaveTaskStorage(serializeTaskStorage(nextTasks)); err != nil {
		return contracts.TaskListResult{}, err
	}

	return buildTaskListResult(nextTasks), nil
}

func (s DefaultTaskListService) RemoveTask(completedTaskText string) (contracts.TaskListResult, error) {
	storageText, err := s.store.LoadTaskStorage()
	if err != nil {
		return contracts.TaskListResult{}, err
	}

	nextTasks := removeTaskByExactText(parseTaskStorage(storageText), completedTaskText)
	if err := s.store.SaveTaskStorage(serializeTaskStorage(nextTasks)); err != nil {
		return contracts.TaskListResult{}, err
	}

	return buildTaskListResult(nextTasks), nil
}

func buildTaskListResult(taskList []string) contracts.TaskListResult {
	return contracts.TaskListResult{
		Tasks: append([]string{}, taskList...),
		Lines: formatTaskList(taskList),
	}
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "8. Green: Add The Task Service"
```
