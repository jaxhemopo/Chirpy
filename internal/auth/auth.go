package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authoraztion header not found")
	}
	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) {
		return "", errors.New("invalid auth header")
	}
	token := strings.TrimPrefix(authHeader, prefix)

	return strings.TrimSpace(token), nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatal(err)
	}
	return match, nil
}

func MakeRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(tokenBytes), nil
}

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}
	const prefix = "ApiKey "
	if len(authHeader) <= len(prefix) {
		return "", errors.New("invalid auth header")
	}
	apiKey := strings.TrimPrefix(authHeader, prefix)

	return strings.TrimSpace(apiKey), nil
}
