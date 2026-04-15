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

func (s DefaultTaskService) GetTask(taskID string) (contracts.Task, error) {
	panic("not implemented")
}

func (s DefaultTaskService) DeleteTask(taskID string) error {
	panic("not implemented")
}
