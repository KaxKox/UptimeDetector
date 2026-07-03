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

func TestValidateToken_Valid(t *testing.T) {
	userID := 42

	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("Nie udało się wygenerować tokenu: %v", err)
	}

	gotUserID, err := ValidateToken(token)
	if err != nil {
		t.Errorf("Oczekiwano braku błędu dla poprawnego tokenu, otrzymano: %v", err)
	}

	if gotUserID != userID {
		t.Errorf("Oczekiwano userID=%d, otrzymano %d", userID, gotUserID)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("to.nie.jest.poprawny.token")

	if err == nil {
		t.Errorf("Oczekiwano błędu dla niepoprawnego tokenu, otrzymano nil")
	}
}

func TestValidateToken_Empty(t *testing.T) {
	_, err := ValidateToken("")

	if err == nil {
		t.Errorf("Oczekiwano błędu dla pustego tokenu, otrzymano nil")
	}
}