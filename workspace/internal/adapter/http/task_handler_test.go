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
