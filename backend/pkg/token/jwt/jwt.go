package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWToken ...
type JWToken struct {
	secretKey       string
	ExpAccessToken  time.Duration
	ExpRefreshToken time.Duration
}

func New(secret string) *JWToken {
	return &JWToken{
		secretKey:       secret,
		ExpAccessToken:  time.Minute * 15,
		ExpRefreshToken: time.Hour * 24,
	}
}

// CreateTokenPair Возвращает пару новых `access` и `refresh` токенов
func (t *JWToken) CreateTokenPair(userID string) (string, string, error) {

	accessToken, err := t.CreateAccessToken(userID) // Create an access token for user id
	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.CreateRefreshToken(userID) // Create a refresh token for user id
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

// CreateAccessToken is a helper function that creates and returns an access token for a given user id
func (t *JWToken) CreateAccessToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // Create a new JWT with HS256 signing method

	token.Claims = jwt.MapClaims{
		"exp":     time.Now().Add(t.ExpAccessToken).Unix(), // Set expiration time to 15 minutes
		"iat":     time.Now().Unix(),                       // Set issued at time to current time
		"type":    "access",                                // Set type to access
		"user_id": userID,                                  // Set user id to given value
	}

	tokenString, err := token.SignedString([]byte(t.secretKey)) // Sign the JWT with secret key and get the string representation
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// CreateRefreshToken is a helper function that creates and returns a refresh token for a given user id
func (t *JWToken) CreateRefreshToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // Create a new JWT with HS256 signing method

	token.Claims = jwt.MapClaims{
		"exp":     time.Now().Add(t.ExpRefreshToken).Unix(), // Set expiration time to 24 hours
		"iat":     time.Now().Unix(),                        // Set issued at time to current time
		"type":    "refresh",                                // Set type to refresh
		"user_id": userID,                                   // Set user id to given value
	}

	tokenString, err := token.SignedString([]byte(t.secretKey)) // Sign the JWT with secret key and get the string representation
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (t *JWToken) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(tt *jwt.Token) (interface{}, error) { // Parse and validate refresh token using secret key
			if _, ok := tt.Method.(*jwt.SigningMethodHMAC); !ok { // Check if signing method is HMAC
				return nil, fmt.Errorf("unexpected signing method: %v", tt.Header["alg"])
			}

			return []byte(t.secretKey), nil // Return secret key as []byte
		},
	)

	if err != nil { // Check if parsing or validation failed
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { // Check if the token is valid
		return claims, err
	}
	return nil, fmt.Errorf("token invalid")
}

func (t *JWToken) RegenerateAccessToken(refreshToken string) (string, error) {
	var err error
	var claims jwt.MapClaims
	claims, err = t.validateToken(refreshToken)
	if err != nil {
		return "", err
	}
	if claims["type"] != "refresh" { // Check if the token type is refresh.
		return "", fmt.Errorf("invalid token type")
	}
	userID := claims["user_id"].(string) // Extract the user id from the claims

	newAccessToken, err := t.CreateAccessToken(userID) // Create a new access token for user id
	if err != nil {
		return "", fmt.Errorf("failed to create access token: %v", err)
	}
	return newAccessToken, nil
}
