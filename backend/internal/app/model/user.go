package model

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// User is a struct that represents a user record in the database
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

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

func (u *User) ComparePasswords(passwd1, passwd2 string) bool {
	if bcrypt.CompareHashAndPassword([]byte(passwd1), []byte(passwd2)) != nil {
		return false
	}
	return true
}
