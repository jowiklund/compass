package session

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type FileStorage string

func (p FileStorage) Save(s *Session) error {
	file, err := os.OpenFile(string(p), os.O_CREATE|os.O_WRONLY, 0600)
	handleErr := func(err error) error {
		return fmt.Errorf("Could not serialize session :: %+v", err)
	}
	if err != nil {
		return handleErr(err)
	}

	jsonData, err := json.Marshal(s)
	if err != nil {
		return handleErr(err)
	}

	if _, err := file.Write(jsonData); err != nil {
		return handleErr(err)
	}
	return nil
}

func (p FileStorage) Load(s *Session) error {
	file, err := os.Open(string(p))
	handleErr := func(err error) error {
		return fmt.Errorf("Could not deserialize session :: %+v", err)
	}
	if err != nil {
		return handleErr(err)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return handleErr(err)
	}

	if err := json.Unmarshal(fileBytes, s); err != nil {
		return handleErr(err)
	}
	return nil
}
