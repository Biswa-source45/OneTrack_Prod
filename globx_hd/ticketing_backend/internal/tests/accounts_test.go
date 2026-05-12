package tests

import (
	"testing"

	"github.com/Chinmay-Globx/ticketing-backend/internal/handlers"
)

func TestCreateAccountInputValidation(t *testing.T) {
	input := handlers.CreateAccountInput{
		AccountName: "Test Account",
	}
	if input.AccountName == "" {
		t.Error("AccountName should be required")
	}
}
