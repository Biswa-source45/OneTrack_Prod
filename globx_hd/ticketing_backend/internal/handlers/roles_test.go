package handlers

import (
	"testing"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
)

func TestIsManagerOrAbove(t *testing.T) {
	role := func(id uint, name string) models.User {
		return models.User{RoleID: id, Role: models.MasterRole{ID: id, RoleName: name}}
	}
	for _, tc := range []struct {
		u    models.User
		want bool
	}{
		{role(2, "manager"), true},
		{role(7, "superadmin"), true},
		{role(9, " SuperAdmin "), true}, // role names in this DB carry stray spaces
		{role(3, "engineer "), false},
		{role(1, "admin"), false},
		{models.User{RoleID: 7}, false}, // role not loaded -> no superadmin rights
	} {
		if got := isManagerOrAbove(tc.u); got != tc.want {
			t.Errorf("RoleID=%d name=%q: got %v want %v", tc.u.RoleID, tc.u.Role.RoleName, got, tc.want)
		}
	}
}
