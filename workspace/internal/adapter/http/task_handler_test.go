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
