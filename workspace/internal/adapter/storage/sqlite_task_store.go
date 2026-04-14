package storageadapter

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	_ "modernc.org/sqlite"
)

type SQLiteTaskStore struct {
	db *sql.DB
}

func NewSQLiteTaskStore(path string) (*SQLiteTaskStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &SQLiteTaskStore{db: db}
	if err := store.ensureSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SQLiteTaskStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteTaskStore) ensureSchema() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL
		)
	`)
	return err
}

func (s *SQLiteTaskStore) ListTasks() ([]contracts.Task, error) {
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

func (s *SQLiteTaskStore) CreateTask(taskText string) (contracts.Task, error) {
	result, err := s.db.Exec(`INSERT INTO tasks (text) VALUES (?)`, taskText)
	if err != nil {
		return contracts.Task{}, err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return contracts.Task{}, err
	}

	return contracts.Task{
		ID:   taskID,
		Text: taskText,
	}, nil
}

func (s *SQLiteTaskStore) GetTask(taskID int64) (contracts.Task, bool, error) {
	var task contracts.Task
	err := s.db.QueryRow(`SELECT id, text FROM tasks WHERE id = ?`, taskID).Scan(&task.ID, &task.Text)
	if err == sql.ErrNoRows {
		return contracts.Task{}, false, nil
	}
	if err != nil {
		return contracts.Task{}, false, err
	}

	return task, true, nil
}

func (s *SQLiteTaskStore) DeleteTask(taskID int64) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
