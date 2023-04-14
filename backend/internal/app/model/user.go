package model

import (
	"fmt"
	"github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

// User is a struct that represents a user record in the database
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;size:100"`
	Username string `json:"username" gorm:"unique;size:100"`
	Password string `json:"password"`
}

// Validate проверяет данные пользователя.
// Длина username от 4 до 100.
// Длина password от 8 до 100.
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Username, validation.Required, validation.Length(4, 100)),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 100)),
	)
}

// EncryptPassword шифрует пароль пользователя
func (u *User) EncryptPassword() error {
	// Hash the password using bcrypt (you need to import "golang.org/x/crypto/bcrypt")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	// Set the hashed password to the user struct
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(passwd2 string) bool {
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd2)) != nil {
		return false
	}
	return true
}
