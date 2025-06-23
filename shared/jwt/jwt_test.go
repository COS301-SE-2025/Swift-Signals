package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken_Success(t *testing.T) {
	secret := []byte("my-secret-key")
	Init(secret)

	userID := "12345"
	role := "admin"
	expiry := time.Minute * 5

	tokenStr, err := GenerateToken(userID, role, expiry)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := ParseToken(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, "swift-signals", claims.Issuer)
}

func TestParseToken_InvalidToken(t *testing.T) {
	secret := []byte("my-secret-key")
	Init(secret)

	invalidToken := "this.is.not.a.jwt"
	_, err := ParseToken(invalidToken)
	assert.Error(t, err)
}

func TestParseToken_WrongSignature(t *testing.T) {
	// Initialize with one secret
	Init([]byte("correct-secret"))
	tokenStr, err := GenerateToken("user1", "user", time.Minute)
	assert.NoError(t, err)

	// Re-initialize with a different secret, simulating wrong signature
	Init([]byte("wrong-secret"))

	_, err = ParseToken(tokenStr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse token")
}

func TestParseToken_Expired(t *testing.T) {
	Init([]byte("secret-key"))

	tokenStr, err := GenerateToken("user1", "user", -time.Minute) // expired token
	assert.NoError(t, err)

	_, err = ParseToken(tokenStr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestInit(t *testing.T) {
	Init([]byte("init-key"))
}
