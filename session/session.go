package session

import (
	"fmt"
)

type Session struct {
	Token     string `json:"token"`
	Validity  string `json:"validity"`
	persister StorageAdaper
}

type StorageAdaper interface {
	Save(*Session) error
	Load(*Session) error
}

type SessionInterface interface {
	Save() error
	Load() error
	GetToken() string
	SetToken(string)
}

// Attempts to open a file with the given sessionPath and write the model as json to it. Will create the file if it doesn't exist.
func (s *Session) Save() error {
	if s.persister == nil {
		return fmt.Errorf("No persister registered, cannot read or write sessions")
	}
	if err := s.persister.Save(s); err != nil {
		return err
	}
	return nil
}

// Attempts to open a file with the given sessionPath and write its' contents to the model.
func (s *Session) Load() error {
	if s.persister == nil {
		return fmt.Errorf("No persister registered, cannot read or write sessions")
	}
	if err := s.persister.Load(s); err != nil {
		return err
	}
	return nil
}

// Mutate the model to set SessionId
// Mutate the model to set Token
func (s *Session) SetToken(id string) {
	s.Token = id
}

// Get token
func (s *Session) GetToken() string {
	return s.Token
}

// Initialize session
func New(options ...func(*Session)) *Session {
	s := &Session{
		persister: nil,
	}
	for _, o := range options {
		o(s)
	}
	return s
}

// Add a storage adapter to handle session storage
func WithStore(p StorageAdaper) func(*Session) {
	return func(s *Session) {
		s.persister = p
	}
}
