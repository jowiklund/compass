package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockResponse struct {
	Message string `json:"message"`
}

func TestClient_Get_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MockResponse{Message: "success"})
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	var result MockResponse
	err := client.Get("/test", &result)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Message != "success" {
		t.Errorf("expected 'success', got '%s'", result.Message)
	}
}

func TestClient_Get_Non200(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	var result MockResponse
	err := client.Get("/test", &result)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestClient_Post_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MockResponse{Message: "posted"})
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	data := map[string]string{"key": "value"}
	var result MockResponse
	err := client.Post("/test", data, &result)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Message != "posted" {
		t.Errorf("expected 'posted', got '%s'", result.Message)
	}
}

func TestClient_Get_WithHeaders(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization header to be 'Bearer test-token'")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL, WithAuth("Bearer test-token"))
	err := client.Get("/test", nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestClient_Get_WithInterceptor(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	interceptorCalled := false

	client := NewClient(mockServer.URL, WithInterceptor(func(c *Client) {
		interceptorCalled = true
	}))

	err := client.Get("/test", nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !interceptorCalled {
		t.Errorf("expected interceptor to be called")
	}
}

func TestClient_Get_WithStatusHandler(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	statusHandlerCalled := false

	client := NewClient(mockServer.URL, WithStatusHandler(http.StatusNotFound, func() {
		statusHandlerCalled = true
	}))

	err := client.Get("/test", nil)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !statusHandlerCalled {
		t.Errorf("expected status handler to be called")
	}
}
