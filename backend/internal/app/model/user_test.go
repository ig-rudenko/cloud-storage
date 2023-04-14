package model_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"web/backend/internal/app/model"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		user    func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			user: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "no username",
			user: func() *model.User {
				u := model.TestUser(t)
				u.Username = ""
				return u
			},
			isValid: false,
		},
		{
			name: "no password",
			user: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "long username",
			user: func() *model.User {
				u := model.TestUser(t)
				longU := make([]rune, 101, 101)
				for i := range longU {
					longU[i] = 'a'
				}
				u.Username = string(longU)
				return u
			},
			isValid: false,
		},
		{
			name: "long password",
			user: func() *model.User {
				u := model.TestUser(t)
				longP := make([]rune, 101, 101)
				for i := range longP {
					longP[i] = 'a'
				}
				u.Password = string(longP)
				return u
			},
			isValid: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.isValid {
				assert.NoError(t, testCase.user().Validate())
			} else {
				assert.Error(t, testCase.user().Validate())
			}
		})
	}
}

func TestUser_ComparePasswords(t *testing.T) {
	user := model.TestUser(t)
	rawPassword := user.Password
	assert.NoError(t, user.EncryptPassword())

	testCases := []struct {
		name    string
		passwd2 string
		isValid bool
	}{
		{
			name:    "valid",
			passwd2: rawPassword,
			isValid: true,
		},
		{
			name:    "with encrypted password input",
			passwd2: user.Password,
			isValid: false,
		},
		{
			name:    "empty password",
			passwd2: "",
			isValid: false,
		},
		{
			name:    "invalid password",
			passwd2: fmt.Sprintf("%s_invalid", rawPassword),
			isValid: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.isValid {
				assert.True(t, user.ComparePassword(testCase.passwd2))
			} else {
				assert.False(t, user.ComparePassword(testCase.passwd2))
			}
		})
	}
}
