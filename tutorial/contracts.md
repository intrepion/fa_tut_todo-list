# Contracts

Create the shared contract file:

```bash
mkdir -p workspace/internal/contracts
touch workspace/internal/contracts/task_api.go
just format
git add --all
git commit --message 'touch workspace/internal/contracts/task_api.go'
```

Put this exact content in `workspace/internal/contracts/task_api.go`:

```go
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
```

Do not add tests here. Keep this layer limited to interfaces, small shared types, and canonical error values.

Then run:

```bash
just format
just check-all
git add --all
git commit --message "Define todo-list REST contracts"
```
