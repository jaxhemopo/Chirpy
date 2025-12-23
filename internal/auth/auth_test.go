package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mysecretkey"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatal("Failed to create JWT:", err)
	}

	validatedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatal("Failed to validate JWT:", err)
	}

	if validatedID != userID {
		t.Fatalf("Expected userID %v, got %v", userID, validatedID)
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "mysecretkey"
	wrongSecret := "mysecretkeys"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatal("Failed to create JWT", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error validation wrong secret")
	}
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "mysecretkey"
	//Can use -time.Hour instead of using time.Sleep to speed up tests.

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatal("Failed to create JWT, err")
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("Expected expired token error")
	}

}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer testing123")
	expected := "testing123"
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatal("Unexpected error getting bearer token:", err)
	}
	if token != expected {
		t.Fatal("Expected token", expected, "got", token)
	}
}
