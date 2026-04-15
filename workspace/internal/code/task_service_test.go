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
