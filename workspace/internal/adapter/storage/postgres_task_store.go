package storageadapter

import (
	"context"
	"database/sql"
	"errors"

	storedb "github.com/intrepion/fa_tut_todo-list/workspace/internal/adapter/storage/db"
	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresTaskStore struct {
	db      *sql.DB
	queries *storedb.Queries
}

func NewPostgresTaskStore(databaseURL string) (*PostgresTaskStore, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	store := &PostgresTaskStore{
		db:      db,
		queries: storedb.New(db),
	}
	if err := store.ensureSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *PostgresTaskStore) Close() error {
	return s.db.Close()
}

func (s *PostgresTaskStore) ensureSchema() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
			text TEXT NOT NULL
		)
	`)
	return err
}

func (s *PostgresTaskStore) ListTasks() ([]contracts.Task, error) {
	rows, err := s.queries.ListTasks(context.Background())
	if err != nil {
		return nil, err
	}

	tasks := make([]contracts.Task, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, contracts.Task{
			ID:   row.ID,
			Text: row.Text,
		})
	}

	return tasks, nil
}

func (s *PostgresTaskStore) CreateTask(taskText string) (contracts.Task, error) {
	row, err := s.queries.CreateTask(context.Background(), taskText)
	if err != nil {
		return contracts.Task{}, err
	}

	return contracts.Task{
		ID:   row.ID,
		Text: row.Text,
	}, nil
}

func (s *PostgresTaskStore) GetTask(taskID string) (contracts.Task, bool, error) {
	row, err := s.queries.GetTask(context.Background(), taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return contracts.Task{}, false, nil
	}
	if err != nil {
		return contracts.Task{}, false, err
	}

	return contracts.Task{
		ID:   row.ID,
		Text: row.Text,
	}, true, nil
}

func (s *PostgresTaskStore) DeleteTask(taskID string) (bool, error) {
	rowsAffected, err := s.queries.DeleteTask(context.Background(), taskID)
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
