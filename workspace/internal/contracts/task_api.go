package contracts

import "errors"

var (
	ErrTaskTextBlank = errors.New("task text must not be blank")
	ErrTaskNotFound  = errors.New("task not found")
)

type Task struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}

type TaskStore interface {
	ListTasks() ([]Task, error)
	CreateTask(taskText string) (Task, error)
	GetTask(taskID int64) (Task, bool, error)
	DeleteTask(taskID int64) (bool, error)
}

type TaskService interface {
	ListTasks() (TaskListResponse, error)
	CreateTask(taskText string) (Task, error)
	GetTask(taskID int64) (Task, error)
	DeleteTask(taskID int64) error
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
