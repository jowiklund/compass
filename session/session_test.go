package session

import "testing"

type MockStorageAdapter struct {
	SavedSession  *Session
	LoadedSession *Session
}

func (m *MockStorageAdapter) Save(session *Session) error {
	m.SavedSession = session
	return nil
}

func (m *MockStorageAdapter) Load(session *Session) error {
	if m.LoadedSession != nil {
		session.Token = m.LoadedSession.Token
	}
	return nil
}

func TestSession_Save(t *testing.T) {
	mockStorage := &MockStorageAdapter{}
	session := New(WithStore(mockStorage))

	session.SetToken("test-token")
	session.Save()

	if mockStorage.SavedSession == nil {
		t.Error("expected session to be saved but it was not")
	}

	if mockStorage.SavedSession.Token != "test-token" {
		t.Errorf("expected Token to be 'test-token', got %s", mockStorage.SavedSession.Token)
	}
}

func TestSession_Load(t *testing.T) {
	mockStorage := &MockStorageAdapter{
		LoadedSession: &Session{Token: "loaded-token"},
	}
	session := New(WithStore(mockStorage))

	session.Load()

	if session.Token != "loaded-token" {
		t.Errorf("expected Token to be 'loaded-token', got %s", session.Token)
	}
}

func TestSession_GetToken_ReturnsTokenWhenSet(t *testing.T) {
	session := &Session{Token: "test-token"}
	credential := session.GetToken()

	if credential != "test-token" {
		t.Errorf("expected credential to be 'test-token', got %s", credential)
	}
}

func TestSession_SetSessionIdAndSetToken(t *testing.T) {
	session := &Session{}
	session.SetToken("new-token")

	if session.Token != "new-token" {
		t.Errorf("expected Token to be 'new-token', got %s", session.Token)
	}
}

func TestNewSession(t *testing.T) {
	session := New()
	if session.persister != nil {
		t.Error("expected persister to be nil")
	}

	mockStorage := &MockStorageAdapter{}
	sessionWithPersister := New(WithStore(mockStorage))
	if sessionWithPersister.persister == nil {
		t.Error("expected persister to be set but it was nil")
	}
}
