package contracts

import "errors"

var (
	ErrTaskTextBlank = errors.New("task text must not be blank")
	ErrTaskNotFound  = errors.New("task not found")
)

type Task struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type TaskStore interface {
	ListTasks() ([]Task, error)
	CreateTask(taskText string) (Task, error)
	GetTask(taskID string) (Task, bool, error)
	DeleteTask(taskID string) (bool, error)
}

type TaskService interface {
	ListTasks() (TaskListResponse, error)
	CreateTask(taskText string) (Task, error)
	GetTask(taskID string) (Task, error)
	DeleteTask(taskID string) error
}

type TaskListResponse struct {
	Tasks []Task `json:"tasks"`
}

type CreateTaskRequest struct {
	Text string `json:"text"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
