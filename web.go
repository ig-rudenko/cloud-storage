package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
)

// User is a struct that represents a user record in the database
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

// TokenPair is a struct that holds an access token and a refresh token
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// DB is a global variable that holds the database connection
var DB *gorm.DB

// SecretKey is a global variable that holds the secret key for signing and validating tokens
var SecretKey = "mysecretkey"

// AuthMiddleware is a function that returns a Gin middleware for checking the authorization token
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization header from the request
		authHeader := c.GetHeader("Authorization")

		// Check if the header is empty or does not start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing or invalid",
			})
			return
		}

		// Extract the token string from the header
		tokenString := authHeader[7:]

		// Parse and validate the token using the secret key
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Check if the signing method is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			// Return the secret key as []byte
			return []byte(secret), nil
		})

		// Check if parsing or validation failed
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("invalid token: %v", err),
			})
			return
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract the user id from the claims
			userID := claims["user_id"].(string)

			// Set the user id as a key-value pair in Gin context
			c.Set("user_id", userID)

			// Proceed to the next handler function
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
	}
}

func main() {

	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/test_go?charset=utf8mb4&parseTime=True&loc=Local"
	// Connect to the database using the username, password, host, port and database name
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the user model
	err = DB.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Create a Gin router with default middleware
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))
	apiRouter := router.Group("/api")

	// Use AuthMiddleware for routes under /api prefix
	apiRouter.Use(AuthMiddleware(SecretKey))

	// Define a route for generating tokens
	router.POST("/token", generateToken)

	router.POST("/register", createUser)

	// Define a route for refreshing tokens
	router.POST("/token/refresh", refreshToken)

	apiRouter.GET("/items/*path", getFilesHandler)
	apiRouter.POST("/items/upload/*path", postFilesHandler)
	apiRouter.POST("/items/create-folder/*path", createFolderHandler)

	apiRouter.POST("/item/rename/*path", renameFileHandler)
	apiRouter.DELETE("/item/*path", deleteFileHandler)

	// Define a route for testing authorization
	apiRouter.GET("/me", func(c *gin.Context) {
		// Get user id from Gin context
		userID, _ := c.Get("user_id")

		// Return user id as JSON response
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
		})
	})

	// Run router on port 8080
	log.Fatal(router.Run(":8080"))
}

func createUser(c *gin.Context) {
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
	if result := DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	// Return the created user as JSON with status code 201
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "id": user.ID})
}

// generateToken is a handler function that creates and returns an access token and a refresh token for a given user credentials
func generateToken(c *gin.Context) {
	var userFormData User
	if err := c.ShouldBindJSON(&userFormData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	result := DB.Where("username = ?", userFormData.Username).First(&user)

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

// refreshToken is a handler function that refreshes and returns an access token for a given refresh token
func refreshToken(c *gin.Context) {
	tokenString := c.PostForm("refresh_token") // Get refresh token from request form

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { // Parse and validate refresh token using secret key
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { // Check if signing method is HMAC
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(SecretKey), nil // Return secret key as []byte
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
			"access_token": newAccessToken,
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

	tokenString, err := token.SignedString([]byte(SecretKey)) // Sign the JWT with secret key and get the string representation
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

	tokenString, err := token.SignedString([]byte(SecretKey)) // Sign the JWT with secret key and get the string representation
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
