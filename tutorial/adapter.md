# Adapter

### 1. Red: Add The List-Tasks Handler Test

Create the first adapter test file:

```bash
touch workspace/internal/adapter/http/task_handler_test.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/http/task_handler_test.go'
```

Put this exact content in `workspace/internal/adapter/http/task_handler_test.go`:

```go
package httpadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskListService struct {
	mock.Mock
}

func (m *MockTaskListService) ListTasks() (contracts.TaskListResult, error) {
	args := m.Called()
	return args.Get(0).(contracts.TaskListResult), args.Error(1)
}

func (m *MockTaskListService) AddTask(taskText string) (contracts.TaskListResult, error) {
	args := m.Called(taskText)
	return args.Get(0).(contracts.TaskListResult), args.Error(1)
}

func (m *MockTaskListService) RemoveTask(completedTaskText string) (contracts.TaskListResult, error) {
	args := m.Called(completedTaskText)
	return args.Get(0).(contracts.TaskListResult), args.Error(1)
}

func TestTaskHandlerGetTasksReturnsTheCurrentTaskList(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskListService)
	service.On("ListTasks").Return(contracts.TaskListResult{
		Tasks: []string{"Learn how to invert binary trees", "Buy milk"},
		Lines: []string{"Learn how to invert binary trees", "Buy milk"},
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.GetTasks(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var body contracts.TaskListResult
	err = json.Unmarshal(rec.Body.Bytes(), &body)

	assert.NoError(t, err)
	assert.Equal(t, []string{"Learn how to invert binary trees", "Buy milk"}, body.Tasks)
	assert.Equal(t, []string{"Learn how to invert binary trees", "Buy milk"}, body.Lines)
	service.AssertExpectations(t)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "1. Red: Add The List-Tasks Handler Test"
```

### 2. Green: Return The Current Task List As JSON

Create the first adapter production file:

```bash
touch workspace/internal/adapter/http/task_handler.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/http/task_handler.go'
```

Put this exact content in `workspace/internal/adapter/http/task_handler.go`:

```go
package httpadapter

import (
	"net/http"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	service contracts.TaskListService
}

func NewTaskHandler(service contracts.TaskListService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
	result, err := h.service.ListTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "2. Green: Return The Current Task List As JSON"
```

### 3. Red: Add The Task-Mutation Handler Tests

Replace `workspace/internal/adapter/http/task_handler_test.go` with:

```go
package httpadapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskListService struct {
	mock.Mock
}

func (m *MockTaskListService) ListTasks() (contracts.TaskListResult, error) {
	args := m.Called()
	return args.Get(0).(contracts.TaskListResult), args.Error(1)
}

func (m *MockTaskListService) AddTask(taskText string) (contracts.TaskListResult, error) {
	args := m.Called(taskText)
	return args.Get(0).(contracts.TaskListResult), args.Error(1)
}

func (m *MockTaskListService) RemoveTask(completedTaskText string) (contracts.TaskListResult, error) {
	args := m.Called(completedTaskText)
	return args.Get(0).(contracts.TaskListResult), args.Error(1)
}

func TestTaskHandlerGetTasksReturnsTheCurrentTaskList(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskListService)
	service.On("ListTasks").Return(contracts.TaskListResult{
		Tasks: []string{"Learn how to invert binary trees", "Buy milk"},
		Lines: []string{"Learn how to invert binary trees", "Buy milk"},
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.GetTasks(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var body contracts.TaskListResult
	err = json.Unmarshal(rec.Body.Bytes(), &body)

	assert.NoError(t, err)
	assert.Equal(t, []string{"Learn how to invert binary trees", "Buy milk"}, body.Tasks)
	assert.Equal(t, []string{"Learn how to invert binary trees", "Buy milk"}, body.Lines)
	service.AssertExpectations(t)
}

func TestTaskHandlerAddTaskAppendsTheSubmittedTask(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/tasks",
		bytes.NewBufferString("{\"task\":\"Buy milk\"}"),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskListService)
	service.On("AddTask", "Buy milk").Return(contracts.TaskListResult{
		Tasks: []string{"Learn how to invert binary trees", "Buy milk"},
		Lines: []string{"Learn how to invert binary trees", "Buy milk"},
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.AddTask(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	service.AssertExpectations(t)
}

func TestTaskHandlerAddTaskRejectsBlankInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/tasks",
		bytes.NewBufferString("{\"task\":\"   \"}"),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskListService)
	handler := NewTaskHandler(service)
	err := handler.AddTask(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertExpectations(t)
}

func TestTaskHandlerRemoveTaskRemovesTheChosenTask(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/tasks?task=Buy%20milk", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskListService)
	service.On("RemoveTask", "Buy milk").Return(contracts.TaskListResult{
		Tasks: []string{"Learn how to invert binary trees"},
		Lines: []string{"Learn how to invert binary trees"},
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.RemoveTask(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	service.AssertExpectations(t)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "3. Red: Add The Task-Mutation Handler Tests"
```

### 4. Green: Handle Add And Remove Task Requests

Replace `workspace/internal/adapter/http/task_handler.go` with:

```go
package httpadapter

import (
	"net/http"
	"strings"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	service contracts.TaskListService
}

func NewTaskHandler(service contracts.TaskListService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
	result, err := h.service.ListTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) AddTask(c echo.Context) error {
	var request contracts.AddTaskRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "invalid request body",
		})
	}

	if strings.TrimSpace(request.Task) == "" {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task must not be blank",
		})
	}

	result, err := h.service.AddTask(request.Task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) RemoveTask(c echo.Context) error {
	taskText := c.QueryParam("task")
	if strings.TrimSpace(taskText) == "" {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task query parameter must not be blank",
		})
	}

	result, err := h.service.RemoveTask(taskText)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "4. Green: Handle Add And Remove Task Requests"
```

### 5. Red: Add The JSON Task Store Test

Create the storage test file:

```bash
touch workspace/internal/adapter/storage/json_task_store_test.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/storage/json_task_store_test.go'
```

Put this exact content in `workspace/internal/adapter/storage/json_task_store_test.go`:

```go
package storageadapter

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONTaskStoreLoadsEmptyJsonArrayWhenTheFileDoesNotExist(t *testing.T) {
	store := NewJSONTaskStore(filepath.Join(t.TempDir(), "tasks.json"))

	result, err := store.LoadTaskStorage()

	require.NoError(t, err)
	assert.JSONEq(t, "[]", result)
}

func TestJSONTaskStoreWritesAndReadsTaskStorage(t *testing.T) {
	store := NewJSONTaskStore(filepath.Join(t.TempDir(), "tasks.json"))

	require.NoError(t, store.SaveTaskStorage("[\"Buy milk\"]"))

	result, err := store.LoadTaskStorage()

	require.NoError(t, err)
	assert.JSONEq(t, "[\"Buy milk\"]", result)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "5. Red: Add The JSON Task Store Test"
```

### 6. Green: Read And Write The Local JSON Task Store

Create the storage production file:

```bash
touch workspace/internal/adapter/storage/json_task_store.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/storage/json_task_store.go'
```

Put this exact content in `workspace/internal/adapter/storage/json_task_store.go`:

```go
package storageadapter

import (
	"errors"
	"os"
	"path/filepath"
)

type JSONTaskStore struct {
	path string
}

func NewJSONTaskStore(path string) *JSONTaskStore {
	return &JSONTaskStore{path: path}
}

func (s *JSONTaskStore) LoadTaskStorage() (string, error) {
	storageBytes, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return "[]", nil
	}
	if err != nil {
		return "", err
	}

	return string(storageBytes), nil
}

func (s *JSONTaskStore) SaveTaskStorage(storageText string) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(s.path, []byte(storageText), 0o644)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "6. Green: Read And Write The Local JSON Task Store"
```

### 7. Green: Wire The Server Entry Point

Create the server entry point:

```bash
touch workspace/cmd/server/main.go
just format
git add --all
git commit --message 'touch workspace/cmd/server/main.go'
```

Put this exact content in `workspace/cmd/server/main.go`:

```go
package main

import (
	"log"

	httpadapter "github.com/intrepion/fa_tut_todo-list/workspace/internal/adapter/http"
	storageadapter "github.com/intrepion/fa_tut_todo-list/workspace/internal/adapter/storage"
	"github.com/intrepion/fa_tut_todo-list/workspace/internal/code"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:25616"},
	}))

	store := storageadapter.NewJSONTaskStore("data/tasks.json")
	service := code.NewTaskListService(store)
	handler := httpadapter.NewTaskHandler(service)

	e.GET("/api/tasks", handler.GetTasks)
	e.POST("/api/tasks", handler.AddTask)
	e.DELETE("/api/tasks", handler.RemoveTask)

	log.Fatal(e.Start(":25664"))
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "7. Green: Wire The Server Entry Point"
```
