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

func (m *MockTaskService) GetTask(taskID string) (contracts.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(taskID string) error {
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
			{ID: "11111111-1111-1111-1111-111111111111", Text: "Learn how to invert binary trees"},
			{ID: "22222222-2222-2222-2222-222222222222", Text: "Buy milk"},
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
		{ID: "11111111-1111-1111-1111-111111111111", Text: "Learn how to invert binary trees"},
		{ID: "22222222-2222-2222-2222-222222222222", Text: "Buy milk"},
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
		ID:   "11111111-1111-1111-1111-111111111111",
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
	taskID := c.Param("id")
	if taskID == "" {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task id must not be empty",
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
	taskID := c.Param("id")
	if taskID == "" {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task id must not be empty",
		})
	}

	err := h.service.DeleteTask(taskID)
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

### 3. Red: Add Fetch, Delete, And Postgres Store Tests

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

func (m *MockTaskService) GetTask(taskID string) (contracts.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(contracts.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(taskID string) error {
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
			{ID: "11111111-1111-1111-1111-111111111111", Text: "Learn how to invert binary trees"},
			{ID: "22222222-2222-2222-2222-222222222222", Text: "Buy milk"},
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
		ID:   "11111111-1111-1111-1111-111111111111",
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
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/11111111-1111-1111-1111-111111111111", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/tasks/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("11111111-1111-1111-1111-111111111111")

	service := new(MockTaskService)
	service.On("GetTask", "11111111-1111-1111-1111-111111111111").Return(contracts.Task{
		ID:   "11111111-1111-1111-1111-111111111111",
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
	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/11111111-1111-1111-1111-111111111111", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/tasks/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("11111111-1111-1111-1111-111111111111")

	service := new(MockTaskService)
	service.On("DeleteTask", "11111111-1111-1111-1111-111111111111").Return(nil)

	handler := NewTaskHandler(service)
	err := handler.DeleteTask(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	service.AssertExpectations(t)
}
```

Create the Postgres store test file:

```bash
mkdir -p workspace/internal/adapter/storage
touch workspace/internal/adapter/storage/postgres_task_store_test.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/storage/postgres_task_store_test.go'
```

Put this exact content in `workspace/internal/adapter/storage/postgres_task_store_test.go`:

```go
package storageadapter

import (
	"regexp"
	"testing"

	storedb "github.com/intrepion/fa_tut_todo-list/workspace/internal/adapter/storage/db"
	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresTaskStoreCreatesTasksWithDatabaseGeneratedUuids(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := &PostgresTaskStore{queries: storedb.New(db)}
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tasks (text) VALUES ($1) RETURNING id, text`)).
		WithArgs("Learn how to invert binary trees").
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("11111111-1111-1111-1111-111111111111", "Learn how to invert binary trees"))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tasks (text) VALUES ($1) RETURNING id, text`)).
		WithArgs("Buy milk").
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("22222222-2222-2222-2222-222222222222", "Buy milk"))

	first, err := store.CreateTask("Learn how to invert binary trees")
	require.NoError(t, err)

	second, err := store.CreateTask("Buy milk")
	require.NoError(t, err)

	assert.Equal(t, contracts.Task{ID: "11111111-1111-1111-1111-111111111111", Text: "Learn how to invert binary trees"}, first)
	assert.Equal(t, contracts.Task{ID: "22222222-2222-2222-2222-222222222222", Text: "Buy milk"}, second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTaskStoreListsGetsAndDeletesTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := &PostgresTaskStore{queries: storedb.New(db)}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, text FROM tasks ORDER BY text ASC, id ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("11111111-1111-1111-1111-111111111111", "Buy milk"))

	listed, err := store.ListTasks()
	require.NoError(t, err)
	assert.Equal(t, []contracts.Task{{ID: "11111111-1111-1111-1111-111111111111", Text: "Buy milk"}}, listed)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, text FROM tasks WHERE id = $1 LIMIT 1`)).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("11111111-1111-1111-1111-111111111111", "Buy milk"))

	got, found, err := store.GetTask("11111111-1111-1111-1111-111111111111")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, contracts.Task{ID: "11111111-1111-1111-1111-111111111111", Text: "Buy milk"}, got)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tasks WHERE id = $1`)).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnResult(sqlmock.NewResult(0, 1))

	deleted, err := store.DeleteTask("11111111-1111-1111-1111-111111111111")
	require.NoError(t, err)
	assert.True(t, deleted)
	assert.NoError(t, mock.ExpectationsWereMet())
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "3. Red: Add Fetch, Delete, And Postgres Store Tests"
```

### 4. Green: Add The Postgres Store And Server Wiring

Use the `sqlc` configuration and SQL files that you created in `tutorial/setup.md`. The generated `just check-tests` and `just run` commands already call `sqlc generate` for you before they compile the app.

Create the Postgres store production file:

```bash
mkdir -p workspace/internal/adapter/storage
touch workspace/internal/adapter/storage/postgres_task_store.go
just format
git add --all
git commit --message 'touch workspace/internal/adapter/storage/postgres_task_store.go'
```

Put this exact content in `workspace/internal/adapter/storage/postgres_task_store.go`:

```go
package storageadapter

import (
	"context"
	"database/sql"
	"errors"

	storedb "github.com/intrepion/fa_tut_todo-list/workspace/internal/adapter/storage/db"
	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresTaskStore struct {
	db      *sql.DB
	queries *storedb.Queries
}

func NewPostgresTaskStore(databaseURL string) (*PostgresTaskStore, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	store := &PostgresTaskStore{
		db:      db,
		queries: storedb.New(db),
	}
	if err := store.ensureSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *PostgresTaskStore) Close() error {
	return s.db.Close()
}

func (s *PostgresTaskStore) ensureSchema() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
			text TEXT NOT NULL
		)
	`)
	return err
}

func (s *PostgresTaskStore) ListTasks() ([]contracts.Task, error) {
	rows, err := s.queries.ListTasks(context.Background())
	if err != nil {
		return nil, err
	}

	tasks := make([]contracts.Task, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, contracts.Task{
			ID:   row.ID,
			Text: row.Text,
		})
	}

	return tasks, nil
}

func (s *PostgresTaskStore) CreateTask(taskText string) (contracts.Task, error) {
	row, err := s.queries.CreateTask(context.Background(), taskText)
	if err != nil {
		return contracts.Task{}, err
	}

	return contracts.Task{
		ID:   row.ID,
		Text: row.Text,
	}, nil
}

func (s *PostgresTaskStore) GetTask(taskID string) (contracts.Task, bool, error) {
	row, err := s.queries.GetTask(context.Background(), taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return contracts.Task{}, false, nil
	}
	if err != nil {
		return contracts.Task{}, false, err
	}

	return contracts.Task{
		ID:   row.ID,
		Text: row.Text,
	}, true, nil
}

func (s *PostgresTaskStore) DeleteTask(taskID string) (bool, error) {
	rowsAffected, err := s.queries.DeleteTask(context.Background(), taskID)
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
	"os"

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

	databaseURL := os.Getenv("TODO_LIST_DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("TODO_LIST_DATABASE_URL must not be empty")
	}

	store, err := storageadapter.NewPostgresTaskStore(databaseURL)
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
git commit --message "4. Green: Add The Postgres Store And Server Wiring"
```
