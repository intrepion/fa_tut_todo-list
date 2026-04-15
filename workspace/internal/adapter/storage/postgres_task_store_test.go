package storageadapter

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresTaskStoreCreatesTasksWithIncreasingIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := &PostgresTaskStore{db: db}
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tasks (text) VALUES ($1) RETURNING id`)).
		WithArgs("Learn how to invert binary trees").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tasks (text) VALUES ($1) RETURNING id`)).
		WithArgs("Buy milk").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(2)))

	first, err := store.CreateTask("Learn how to invert binary trees")
	require.NoError(t, err)

	second, err := store.CreateTask("Buy milk")
	require.NoError(t, err)

	assert.Equal(t, contracts.Task{ID: 1, Text: "Learn how to invert binary trees"}, first)
	assert.Equal(t, contracts.Task{ID: 2, Text: "Buy milk"}, second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTaskStoreListsGetsAndDeletesTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := &PostgresTaskStore{db: db}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, text FROM tasks ORDER BY id ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow(int64(1), "Buy milk"))

	listed, err := store.ListTasks()
	require.NoError(t, err)
	assert.Equal(t, []contracts.Task{{ID: 1, Text: "Buy milk"}}, listed)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, text FROM tasks WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow(int64(1), "Buy milk"))

	got, found, err := store.GetTask(1)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, contracts.Task{ID: 1, Text: "Buy milk"}, got)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tasks WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	deleted, err := store.DeleteTask(1)
	require.NoError(t, err)
	assert.True(t, deleted)
	assert.NoError(t, mock.ExpectationsWereMet())
}
