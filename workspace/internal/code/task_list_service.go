package code

import "encoding/json"

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
