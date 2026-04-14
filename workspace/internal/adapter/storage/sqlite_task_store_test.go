package storageadapter

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteTaskStoreCreatesTasksWithIncreasingIds(t *testing.T) {
	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "tasks.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = store.Close() })

	first, err := store.CreateTask("Learn how to invert binary trees")
	require.NoError(t, err)

	second, err := store.CreateTask("Buy milk")
	require.NoError(t, err)

	assert.Equal(t, int64(1), first.ID)
	assert.Equal(t, int64(2), second.ID)
}

func TestSQLiteTaskStoreListsGetsAndDeletesTasks(t *testing.T) {
	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "tasks.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = store.Close() })

	created, err := store.CreateTask("Buy milk")
	require.NoError(t, err)

	listed, err := store.ListTasks()
	require.NoError(t, err)
	assert.Len(t, listed, 1)

	got, found, err := store.GetTask(created.ID)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, created, got)

	deleted, err := store.DeleteTask(created.ID)
	require.NoError(t, err)
	assert.True(t, deleted)
}
