package storage

import (
	"encoding/json"
	"io"
	"to-do-list/internal/tasks"
)

type Storage struct {
	file io.ReadWriter
}
func NewStorage(file io.ReadWriter) *Storage {
    return &Storage{
        file: file,
    }
}


func (s *Storage) Create(task tasks.Task) error {
	bytes, err := json.Marshal(&task)
	if err != nil {
		return err
	}

	_, err = s.file.Read(bytes)
	if err != nil {
		return err
	}
	return nil
}
