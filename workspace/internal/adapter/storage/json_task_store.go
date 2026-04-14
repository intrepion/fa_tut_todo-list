package storageadapter

import (
	"errors"
	"os"
	"path/filepath"
)

type JSONTaskStore struct {
	path string
}

func NewJSONTaskStore(path string) *JSONTaskStore {
	return &JSONTaskStore{path: path}
}

func (s *JSONTaskStore) LoadTaskStorage() (string, error) {
	storageBytes, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return "[]", nil
	}
	if err != nil {
		return "", err
	}

	return string(storageBytes), nil
}

func (s *JSONTaskStore) SaveTaskStorage(storageText string) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(s.path, []byte(storageText), 0o644)
}
