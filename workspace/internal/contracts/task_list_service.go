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
