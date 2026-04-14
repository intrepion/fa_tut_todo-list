package code

import (
	"encoding/json"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
)

type DefaultTaskListService struct {
	store contracts.TaskStore
}

func NewTaskListService(store contracts.TaskStore) contracts.TaskListService {
	return DefaultTaskListService{store: store}
}

func parseTaskStorage(storageText string) []string {
	var tasks []string
	_ = json.Unmarshal([]byte(storageText), &tasks)
	return tasks
}

func appendTask(taskList []string, taskText string) []string {
	next := append([]string{}, taskList...)
	return append(next, taskText)
}

func removeTaskByExactText(taskList []string, completedTaskText string) []string {
	next := append([]string{}, taskList...)
	for index, task := range next {
		if task == completedTaskText {
			return append(next[:index], next[index+1:]...)
		}
	}
	return next
}

func formatTaskList(taskList []string) []string {
	return append([]string{}, taskList...)
}

func serializeTaskStorage(taskList []string) string {
	storageBytes, _ := json.Marshal(taskList)
	return string(storageBytes)
}

func (s DefaultTaskListService) ListTasks() (contracts.TaskListResult, error) {
	storageText, err := s.store.LoadTaskStorage()
	if err != nil {
		return contracts.TaskListResult{}, err
	}

	tasks := parseTaskStorage(storageText)
	return buildTaskListResult(tasks), nil
}

func (s DefaultTaskListService) AddTask(taskText string) (contracts.TaskListResult, error) {
	storageText, err := s.store.LoadTaskStorage()
	if err != nil {
		return contracts.TaskListResult{}, err
	}

	nextTasks := appendTask(parseTaskStorage(storageText), taskText)
	if err := s.store.SaveTaskStorage(serializeTaskStorage(nextTasks)); err != nil {
		return contracts.TaskListResult{}, err
	}

	return buildTaskListResult(nextTasks), nil
}

func (s DefaultTaskListService) RemoveTask(completedTaskText string) (contracts.TaskListResult, error) {
	storageText, err := s.store.LoadTaskStorage()
	if err != nil {
		return contracts.TaskListResult{}, err
	}

	nextTasks := removeTaskByExactText(parseTaskStorage(storageText), completedTaskText)
	if err := s.store.SaveTaskStorage(serializeTaskStorage(nextTasks)); err != nil {
		return contracts.TaskListResult{}, err
	}

	return buildTaskListResult(nextTasks), nil
}

func buildTaskListResult(taskList []string) contracts.TaskListResult {
	return contracts.TaskListResult{
		Tasks: append([]string{}, taskList...),
		Lines: formatTaskList(taskList),
	}
}
