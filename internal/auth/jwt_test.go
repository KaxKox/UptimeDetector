package auth

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	userID := 1

	token, err := GenerateToken(userID)

	if err != nil {
		t.Errorf("Oczekiwano braku błędu, otrzymano błąd: %v", err)
	}

	if token == "" {
		t.Errorf("Oczekiwano tokenu, otrzymano pusty tekst")
	}
}