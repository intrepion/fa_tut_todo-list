package storageadapter

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	storedb "github.com/intrepion/fa_tut_todo-list/workspace/internal/adapter/storage/db"
	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresTaskStoreCreatesTasksWithDatabaseGeneratedUuids(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := &PostgresTaskStore{queries: storedb.New(db)}
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tasks (text) VALUES ($1) RETURNING id, text`)).
		WithArgs("Learn how to invert binary trees").
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("11111111-1111-1111-1111-111111111111", "Learn how to invert binary trees"))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tasks (text) VALUES ($1) RETURNING id, text`)).
		WithArgs("Buy milk").
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("22222222-2222-2222-2222-222222222222", "Buy milk"))

	first, err := store.CreateTask("Learn how to invert binary trees")
	require.NoError(t, err)

	second, err := store.CreateTask("Buy milk")
	require.NoError(t, err)

	assert.Equal(t, contracts.Task{ID: "11111111-1111-1111-1111-111111111111", Text: "Learn how to invert binary trees"}, first)
	assert.Equal(t, contracts.Task{ID: "22222222-2222-2222-2222-222222222222", Text: "Buy milk"}, second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTaskStoreListsGetsAndDeletesTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := &PostgresTaskStore{queries: storedb.New(db)}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, text FROM tasks ORDER BY text ASC, id ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("11111111-1111-1111-1111-111111111111", "Buy milk"))

	listed, err := store.ListTasks()
	require.NoError(t, err)
	assert.Equal(t, []contracts.Task{{ID: "11111111-1111-1111-1111-111111111111", Text: "Buy milk"}}, listed)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, text FROM tasks WHERE id = $1 LIMIT 1`)).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("11111111-1111-1111-1111-111111111111", "Buy milk"))

	got, found, err := store.GetTask("11111111-1111-1111-1111-111111111111")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, contracts.Task{ID: "11111111-1111-1111-1111-111111111111", Text: "Buy milk"}, got)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tasks WHERE id = $1`)).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnResult(sqlmock.NewResult(0, 1))

	deleted, err := store.DeleteTask("11111111-1111-1111-1111-111111111111")
	require.NoError(t, err)
	assert.True(t, deleted)
	assert.NoError(t, mock.ExpectationsWereMet())
}
