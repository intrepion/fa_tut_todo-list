# Adapter

### 1. Red: Add The REST Handler Tests

Create the first adapter test file:

```bash
mkdir -p workspace/internal/adapter/http
touch workspace/internal/adapter/http/task_handler_test.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/http/task_handler_test.go'
```

Put this exact content in `workspace/internal/adapter/http/task_handler_test.go`:

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

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) ListTasks() (contracts.TaskListResponse, error) {
	args := m.Called()
	return args.Get(0).(contracts.TaskListResponse), args.Error(1)
}

func (m *MockTaskService) CreateTask(taskText string) (contracts.Task, error) {
	args := m.Called(taskText)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskService) GetTask(taskID int64) (contracts.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(taskID int64) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func TestTaskHandlerListTasksReturnsTheCurrentTaskResources(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskService)
	service.On("ListTasks").Return(contracts.TaskListResponse{
		Tasks: []contracts.Task{
			{ID: 1, Text: "Learn how to invert binary trees"},
			{ID: 2, Text: "Buy milk"},
		},
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.ListTasks(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var body contracts.TaskListResponse
	err = json.Unmarshal(rec.Body.Bytes(), &body)

	assert.NoError(t, err)
	assert.Equal(t, []contracts.Task{
		{ID: 1, Text: "Learn how to invert binary trees"},
		{ID: 2, Text: "Buy milk"},
	}, body.Tasks)
	service.AssertExpectations(t)
}

func TestTaskHandlerCreateTaskReturnsCreatedResource(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/tasks",
		bytes.NewBufferString("{\"text\":\"Buy milk\"}"),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskService)
	service.On("CreateTask", "Buy milk").Return(contracts.Task{
		ID:   1,
		Text: "Buy milk",
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.CreateTask(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	service.AssertExpectations(t)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "1. Red: Add The REST Handler Tests"
```

### 2. Green: Return REST Resources

Create the first adapter production file:

```bash
mkdir -p workspace/internal/adapter/http
touch workspace/internal/adapter/http/task_handler.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/http/task_handler.go'
```

Put this exact content in `workspace/internal/adapter/http/task_handler.go`:

```go
package httpadapter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	service contracts.TaskService
}

func NewTaskHandler(service contracts.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) ListTasks(c echo.Context) error {
	result, err := h.service.ListTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
	var request contracts.CreateTaskRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "invalid request body",
		})
	}

	task, err := h.service.CreateTask(request.Text)
	if err != nil {
		if errors.Is(err, contracts.ErrTaskTextBlank) {
			return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task id must be an integer",
		})
	}

	task, err := h.service.GetTask(taskID)
	if err != nil {
		if errors.Is(err, contracts.ErrTaskNotFound) {
			return c.JSON(http.StatusNotFound, contracts.ErrorResponse{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task id must be an integer",
		})
	}

	err = h.service.DeleteTask(taskID)
	if err != nil {
		if errors.Is(err, contracts.ErrTaskNotFound) {
			return c.JSON(http.StatusNotFound, contracts.ErrorResponse{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "2. Green: Return REST Resources"
```

### 3. Red: Add Fetch, Delete, And SQLite Store Tests

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

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) ListTasks() (contracts.TaskListResponse, error) {
	args := m.Called()
	return args.Get(0).(contracts.TaskListResponse), args.Error(1)
}

func (m *MockTaskService) CreateTask(taskText string) (contracts.Task, error) {
	args := m.Called(taskText)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskService) GetTask(taskID int64) (contracts.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(taskID int64) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func TestTaskHandlerListTasksReturnsTheCurrentTaskResources(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskService)
	service.On("ListTasks").Return(contracts.TaskListResponse{
		Tasks: []contracts.Task{
			{ID: 1, Text: "Learn how to invert binary trees"},
			{ID: 2, Text: "Buy milk"},
		},
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.ListTasks(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	service.AssertExpectations(t)
}

func TestTaskHandlerCreateTaskReturnsCreatedResource(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/tasks",
		bytes.NewBufferString("{\"text\":\"Buy milk\"}"),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	service := new(MockTaskService)
	service.On("CreateTask", "Buy milk").Return(contracts.Task{
		ID:   1,
		Text: "Buy milk",
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.CreateTask(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	service.AssertExpectations(t)
}

func TestTaskHandlerGetTaskReturnsTheRequestedResource(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/tasks/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	service := new(MockTaskService)
	service.On("GetTask", int64(1)).Return(contracts.Task{
		ID:   1,
		Text: "Buy milk",
	}, nil)

	handler := NewTaskHandler(service)
	err := handler.GetTask(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	service.AssertExpectations(t)
}

func TestTaskHandlerDeleteTaskReturnsNoContent(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/1", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/tasks/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	service := new(MockTaskService)
	service.On("DeleteTask", int64(1)).Return(nil)

	handler := NewTaskHandler(service)
	err := handler.DeleteTask(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	service.AssertExpectations(t)
}
```

Create the SQLite store test file:

```bash
mkdir -p workspace/internal/adapter/storage
touch workspace/internal/adapter/storage/sqlite_task_store_test.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/storage/sqlite_task_store_test.go'
```

Put this exact content in `workspace/internal/adapter/storage/sqlite_task_store_test.go`:

```go
package storageadapter

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteTaskStoreCreatesTasksWithIncreasingIds(t *testing.T) {
	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "tasks.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = store.Close() })

	first, err := store.CreateTask("Learn how to invert binary trees")
	require.NoError(t, err)

	second, err := store.CreateTask("Buy milk")
	require.NoError(t, err)

	assert.Equal(t, int64(1), first.ID)
	assert.Equal(t, int64(2), second.ID)
}

func TestSQLiteTaskStoreListsGetsAndDeletesTasks(t *testing.T) {
	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "tasks.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = store.Close() })

	created, err := store.CreateTask("Buy milk")
	require.NoError(t, err)

	listed, err := store.ListTasks()
	require.NoError(t, err)
	assert.Len(t, listed, 1)

	got, found, err := store.GetTask(created.ID)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, created, got)

	deleted, err := store.DeleteTask(created.ID)
	require.NoError(t, err)
	assert.True(t, deleted)
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "3. Red: Add Fetch, Delete, And SQLite Store Tests"
```

### 4. Green: Add The SQLite Store And Server Wiring

Create the SQLite store production file:

```bash
mkdir -p workspace/internal/adapter/storage
touch workspace/internal/adapter/storage/sqlite_task_store.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/storage/sqlite_task_store.go'
```

Put this exact content in `workspace/internal/adapter/storage/sqlite_task_store.go`:

```go
package storageadapter

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	_ "modernc.org/sqlite"
)

type SQLiteTaskStore struct {
	db *sql.DB
}

func NewSQLiteTaskStore(path string) (*SQLiteTaskStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &SQLiteTaskStore{db: db}
	if err := store.ensureSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SQLiteTaskStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteTaskStore) ensureSchema() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL
		)
	`)
	return err
}

func (s *SQLiteTaskStore) ListTasks() ([]contracts.Task, error) {
	rows, err := s.db.Query(`SELECT id, text FROM tasks ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []contracts.Task
	for rows.Next() {
		var task contracts.Task
		if err := rows.Scan(&task.ID, &task.Text); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (s *SQLiteTaskStore) CreateTask(taskText string) (contracts.Task, error) {
	result, err := s.db.Exec(`INSERT INTO tasks (text) VALUES (?)`, taskText)
	if err != nil {
		return contracts.Task{}, err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return contracts.Task{}, err
	}

	return contracts.Task{
		ID:   taskID,
		Text: taskText,
	}, nil
}

func (s *SQLiteTaskStore) GetTask(taskID int64) (contracts.Task, bool, error) {
	var task contracts.Task
	err := s.db.QueryRow(`SELECT id, text FROM tasks WHERE id = ?`, taskID).Scan(&task.ID, &task.Text)
	if err == sql.ErrNoRows {
		return contracts.Task{}, false, nil
	}
	if err != nil {
		return contracts.Task{}, false, err
	}

	return task, true, nil
}

func (s *SQLiteTaskStore) DeleteTask(taskID int64) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
```

Create the server entry point:

```bash
mkdir -p workspace/cmd/server
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

	store, err := storageadapter.NewSQLiteTaskStore("data/tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	service := code.NewTaskService(store)
	handler := httpadapter.NewTaskHandler(service)

	e.GET("/api/tasks", handler.ListTasks)
	e.POST("/api/tasks", handler.CreateTask)
	e.GET("/api/tasks/:id", handler.GetTask)
	e.DELETE("/api/tasks/:id", handler.DeleteTask)

	log.Fatal(e.Start(":25664"))
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "4. Green: Add The SQLite Store And Server Wiring"
```
