package util

import (
	"errors"
	"net/http"
	"strings"
)

func GetToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	token := strings.TrimSpace(authHeader[len(prefix):])
	if token == "" {
		return "", errors.New("token is missing after 'Bearer '")
	}

	return token, nil
}
