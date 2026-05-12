package tests

import (
	"testing"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
)

func TestAccountModel(t *testing.T) {
	acc := models.Account{
		AccountName:  "Test",
		CustomerCode: "001",
	}
	if acc.AccountName == "" {
		t.Error("AccountName should not be empty")
	}
	if acc.CustomerCode == "" {
		t.Error("CustomerCode should not be empty")
	}
}
