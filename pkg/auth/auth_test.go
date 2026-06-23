package auth

import (
	"testing"
	"time"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "my-secret-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned empty string")
	}

	if hash == password {
		t.Fatal("HashPassword returned plain text")
	}

	if !CheckPassword(password, hash) {
		t.Fatal("CheckPassword returned false for correct password")
	}

	if CheckPassword("wrong-password", hash) {
		t.Fatal("CheckPassword returned true for incorrect password")
	}
}

func TestGenerateTokenAndValidateToken(t *testing.T) {
	secret := "this-is-a-very-long-secret-key-for-jwt-signing"
	expiration := time.Hour

	token, err := GenerateToken(secret, expiration, "user-123", "test@example.com", "owner")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("GenerateToken returned empty string")
	}

	claims, err := ValidateToken(secret, token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user-123")
	}

	if claims.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", claims.Email, "test@example.com")
	}

	if claims.Role != "owner" {
		t.Errorf("Role = %q, want %q", claims.Role, "owner")
	}
}

func TestValidateTokenExpired(t *testing.T) {
	secret := "this-is-a-very-long-secret-key-for-jwt-signing"

	token, err := GenerateToken(secret, -time.Hour, "user-123", "test@example.com", "owner")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ValidateToken(secret, token)
	if err == nil {
		t.Fatal("ValidateToken did not return error for expired token")
	}
}

func TestValidateTokenInvalidSecret(t *testing.T) {
	secret := "this-is-a-very-long-secret-key-for-jwt-signing"

	token, err := GenerateToken(secret, time.Hour, "user-123", "test@example.com", "owner")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ValidateToken("different-secret-key-for-jwt-signing-tests", token)
	if err == nil {
		t.Fatal("ValidateToken did not return error for invalid secret")
	}
}
