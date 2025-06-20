package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	userpb "github.com/COS301-SE-2025/Swift-Signals/api-gateway/protos/user"
)

func TestPostLogin_Success(t *testing.T) {
	mockClient := &MockUserClient{
		LoginUserFunc: func(ctx context.Context, email, password string) (*userpb.AuthResponse, error) {
			return &userpb.AuthResponse{Token: "test-token"}, nil
		},
	}

	handler := &AuthHandler{UserClient: mockClient}

	body := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	handler.PostLogin(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if loginResp.Token != "test-token" {
		t.Errorf("expected token 'test-token', got '%s'", loginResp.Token)
	}
}
