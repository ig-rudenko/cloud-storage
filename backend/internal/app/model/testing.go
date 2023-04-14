package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		ID:       0,
		Username: "test_user",
		Password: "password",
	}
}
