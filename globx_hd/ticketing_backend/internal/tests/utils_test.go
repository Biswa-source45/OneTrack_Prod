package tests

import (
	"testing"

	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
)

func TestHashPasswordAndCheck(t *testing.T) {
	pw := "password123"
	hash, err := utils.HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if !utils.CheckPasswordHash(pw, hash) {
		t.Error("CheckPasswordHash should return true for correct password")
	}
	if utils.CheckPasswordHash("wrongpw", hash) {
		t.Error("CheckPasswordHash should return false for wrong password")
	}
}
