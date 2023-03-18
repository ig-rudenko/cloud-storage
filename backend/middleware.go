package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

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
