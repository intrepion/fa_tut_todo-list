package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
