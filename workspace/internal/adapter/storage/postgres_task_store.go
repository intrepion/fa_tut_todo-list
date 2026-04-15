package storageadapter

import (
	"database/sql"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(databaseURL string) (*PostgresTaskStore, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	store := &PostgresTaskStore{db: db}
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
			id BIGSERIAL PRIMARY KEY,
			text TEXT NOT NULL
		)
	`)
	return err
}

func (s *PostgresTaskStore) ListTasks() ([]contracts.Task, error) {
	rows, err := s.db.Query(`SELECT id, text FROM tasks ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []contracts.Task
	for rows.Next() {
		var task contracts.Task
		if err := rows.Scan(&task.ID, &task.Text); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (s *PostgresTaskStore) CreateTask(taskText string) (contracts.Task, error) {
	var taskID int64
	if err := s.db.QueryRow(
		`INSERT INTO tasks (text) VALUES ($1) RETURNING id`,
		taskText,
	).Scan(&taskID); err != nil {
		return contracts.Task{}, err
	}

	return contracts.Task{
		ID:   taskID,
		Text: taskText,
	}, nil
}

func (s *PostgresTaskStore) GetTask(taskID int64) (contracts.Task, bool, error) {
	var task contracts.Task
	err := s.db.QueryRow(
		`SELECT id, text FROM tasks WHERE id = $1`,
		taskID,
	).Scan(&task.ID, &task.Text)
	if err == sql.ErrNoRows {
		return contracts.Task{}, false, nil
	}
	if err != nil {
		return contracts.Task{}, false, err
	}

	return task, true, nil
}

func (s *PostgresTaskStore) DeleteTask(taskID int64) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM tasks WHERE id = $1`, taskID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
