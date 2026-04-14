# Contracts

Create the shared contract file:

```bash
touch workspace/internal/contracts/task_list_service.go
just format
just check-all
git add --all
git commit --message 'touch workspace/internal/contracts/task_list_service.go'
```

Put this exact content in `workspace/internal/contracts/task_list_service.go`:

```go
package contracts

type TaskStore interface {
	LoadTaskStorage() (string, error)
	SaveTaskStorage(storageText string) error
}

type TaskListService interface {
	ListTasks() (TaskListResult, error)
	AddTask(taskText string) (TaskListResult, error)
	RemoveTask(completedTaskText string) (TaskListResult, error)
}

type TaskListResult struct {
	Tasks []string `json:"tasks"`
	Lines []string `json:"lines"`
}

type AddTaskRequest struct {
	Task string `json:"task"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
```

Do not add tests here. Keep this layer limited to interfaces and small shared types.

Then run:

```bash
just format
just check-all
git add --all
git commit --message "Define task-list contracts"
```
