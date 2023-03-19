package api

// TokenPair is a struct that holds an access token and a refresh token
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// User is a struct that represents a user record in the database
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

type NewFileName struct {
	Name string `json:"newName"`
}
