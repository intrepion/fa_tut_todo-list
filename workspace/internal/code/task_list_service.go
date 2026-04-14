package code

import "encoding/json"

func parseTaskStorage(storageText string) []string {
	var tasks []string
	_ = json.Unmarshal([]byte(storageText), &tasks)
	return tasks
}
