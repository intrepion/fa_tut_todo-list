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
