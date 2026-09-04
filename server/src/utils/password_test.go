package utils

import (
	"testing"
)

func TestHashPasswordAndCheck(t *testing.T) {
	password := "SecretPass123!"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hashed) {
		t.Errorf("Expected password hash check to succeed")
	}

	if CheckPasswordHash("WrongPassword", hashed) {
		t.Errorf("Expected wrong password check to fail")
	}
}
