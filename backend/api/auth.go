package api

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
	"web/backend/config"
)

// handleDBEntryError Обрабатываем ошибку базы данных. Определяем текст пользователю и статус ответа
func handleDBEntryError(err error) (string, int) {
	var mysqlErr *mysql.MySQLError           // define a variable of type *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok { // check if err can be assigned to mysqlErr variable (type assertion)
		if mysqlErr.Number == 1062 { // error number for duplicate entry https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html#error_er_dup_entry
			return "A user already exists with the same username", http.StatusBadRequest // return from the function or do something else
		}
	}
	return err.Error(), http.StatusInternalServerError // otherwise print the original error
}

// CreateUser Регистрация нового пользователя
func CreateUser(c *gin.Context) {
	// Bind the JSON data to a user struct
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password using bcrypt (you need to import "golang.org/x/crypto/bcrypt")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Set the hashed password to the user struct
	user.Password = string(hashedPassword)

	// Save the user to the database
	if err := config.DB.Create(&user).Error; err != nil {
		errorText, statusCode := handleDBEntryError(err) // Обрабатываем ошибку создания нового пользователя
		c.JSON(statusCode, gin.H{"error": errorText})
		return
	}

	err = os.MkdirAll(config.StorageDir+"/"+strconv.Itoa(int(user.ID)), 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created user as JSON with status code 201
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "id": user.ID})
}

// GenerateToken is a handler function that creates and returns an access token and a refresh token for a given user credentials
func GenerateToken(c *gin.Context) {
	var userFormData User
	if err := c.ShouldBindJSON(&userFormData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	result := config.DB.Where("username = ?", userFormData.Username).First(&user)

	// Check for errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid username"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userFormData.Password)) != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	tokenPair, err := createTokenPair(strconv.Itoa(int(user.ID))) // Create a token pair for user id
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to create tokens: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, tokenPair) // Return token pair as JSON with 200 status code

}

// RefreshToken is a handler function that refreshes and returns an access token for a given refresh token
func RefreshToken(c *gin.Context) {

	var t TokenPair
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(t.RefreshToken, func(t *jwt.Token) (interface{}, error) { // Parse and validate refresh token using secret key
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { // Check if signing method is HMAC
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.SecretKey), nil // Return secret key as []byte
	})

	if err != nil { // Check if parsing or validation failed
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("invalid refresh token: %v", err),
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { // Check if the token is valid
		if claims["type"] != "refresh" { // Check if the token type is refresh.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token type",
			})
			return
		}

		userID := claims["user_id"].(string) // Extract the user id from the claims

		newAccessToken, err := createAccessToken(userID) // Create a new access token for user id
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to create access token: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{ // Return new access token as JSON with 200 status code
			"accessToken": newAccessToken,
		})

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}
}

// createTokenPair is a helper function that creates and returns an access token and a refresh token for a given user id
func createTokenPair(userID string) (TokenPair, error) {
	var tokenPair TokenPair

	accessToken, err := createAccessToken(userID) // Create an access token for user id
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := createRefreshToken(userID) // Create a refresh token for user id
	if err != nil {
		return TokenPair{}, err
	}

	// Set the access and refresh tokens in the TokenPair struct
	tokenPair.AccessToken = accessToken
	tokenPair.RefreshToken = refreshToken

	return tokenPair, nil

}

// createAccessToken is a helper function that creates and returns an access token for a given user id
func createAccessToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // Create a new JWT with HS256 signing method

	token.Claims = jwt.MapClaims{ // Set some standard and custom claims
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Set expiration time to 15 minutes
		"iat":     time.Now().Unix(),                       // Set issued at time to current time
		"type":    "access",                                // Set type to access
		"user_id": userID,                                  // Set user id to given value
	}

	tokenString, err := token.SignedString([]byte(config.SecretKey)) // Sign the JWT with secret key and get the string representation
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// createRefreshToken is a helper function that creates and returns a refresh token for a given user id
func createRefreshToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // Create a new JWT with HS256 signing method

	token.Claims = jwt.MapClaims{ // Set some standard and custom claims
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Set expiration time to 24 hours
		"iat":     time.Now().Unix(),                     // Set issued at time to current time
		"type":    "refresh",                             // Set type to refresh
		"user_id": userID,                                // Set user id to given value
	}

	tokenString, err := token.SignedString([]byte(config.SecretKey)) // Sign the JWT with secret key and get the string representation
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
