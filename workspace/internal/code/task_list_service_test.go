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
