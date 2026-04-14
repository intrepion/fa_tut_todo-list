package storageadapter

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONTaskStoreLoadsEmptyJsonArrayWhenTheFileDoesNotExist(t *testing.T) {
	store := NewJSONTaskStore(filepath.Join(t.TempDir(), "tasks.json"))

	result, err := store.LoadTaskStorage()

	require.NoError(t, err)
	assert.JSONEq(t, "[]", result)
}

func TestJSONTaskStoreWritesAndReadsTaskStorage(t *testing.T) {
	store := NewJSONTaskStore(filepath.Join(t.TempDir(), "tasks.json"))

	require.NoError(t, store.SaveTaskStorage("[\"Buy milk\"]"))

	result, err := store.LoadTaskStorage()

	require.NoError(t, err)
	assert.JSONEq(t, "[\"Buy milk\"]", result)
}
