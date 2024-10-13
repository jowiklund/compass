package session

import (
	"encoding/json"
	"os"
	"testing"
)

func newTestSession() *Session {
	return &Session{
		Token: "test-token",
	}
}

func TestFileStorage_Save(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "session_test_*.json")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	storage := FileStorage(tmpFile.Name())
	session := newTestSession()

	storage.Save(session)

	fileBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("unable to read temp file: %v", err)
	}

	var savedSession Session
	if err := json.Unmarshal(fileBytes, &savedSession); err != nil {
		t.Fatalf("unable to unmarshal session from file: %v", err)
	}

	if savedSession.Token != session.Token {
		t.Errorf("expected Token to be %s, got %s", session.Token, savedSession.Token)
	}
}

func TestFileStorage_Load(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "session_test_*.json")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	session := newTestSession()
	jsonData, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("unable to marshal session: %v", err)
	}
	if err := os.WriteFile(tmpFile.Name(), jsonData, 0600); err != nil {
		t.Fatalf("unable to write session to temp file: %v", err)
	}

	storage := FileStorage(tmpFile.Name())
	loadedSession := &Session{}
	storage.Load(loadedSession)

	if loadedSession.Token != session.Token {
		t.Errorf("expected Token to be %s, got %s", session.Token, loadedSession.Token)
	}
}

func TestFileStorage_Load_NonExistentFile(t *testing.T) {
	nonExistentFile := "non-existent-file.json"
	storage := FileStorage(nonExistentFile)

	session := &Session{}
	storage.Load(session)

	if session.Token != "" {
		t.Errorf("expected an empty session, but got %+v", session)
	}
}

func TestFileStorage_Save_InvalidPath(t *testing.T) {
	invalidPath := "/invalid-path/session.json"
	storage := FileStorage(invalidPath)
	session := newTestSession()

	storage.Save(session)

	// Make sure the program doesn't crash and logs an error
}
