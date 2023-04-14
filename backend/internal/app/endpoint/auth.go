package endpoint

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web/backend/internal/app/model"
)

// RegisterNewUser Регистрация нового пользователя
func (e *Endpoint) RegisterNewUser(c *gin.Context) {
	// Bind the JSON data to a user struct
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем данные пользователя
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем пользователя, а также инициализируем новое пользовательское хранилище
	if err := e.service.InitUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the created user as JSON with status code 201
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "id": user.ID})
}

// GenerateToken is a handler function that creates and returns an access token and a refresh token for a given user credentials
func (e *Endpoint) GenerateToken(c *gin.Context) {
	var userFormData model.User
	if err := c.ShouldBindJSON(&userFormData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := e.service.GetUser(userFormData.Username)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid username"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if user.ComparePassword(userFormData.Password) == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	userID := strconv.Itoa(int(user.ID))                                 // переводим ID пользователя в строку
	accessToken, refreshToken, err := e.tokenGen.CreateTokenPair(userID) // Create a token pair for user id
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to create tokens: %v", err),
		})
		return
	}

	// Return token pair as JSON with 200 status code
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})

}

// RefreshToken is a handler function that refreshes and returns an access token for a given refresh token
func (e *Endpoint) RefreshToken(c *gin.Context) {
	var token TokenPair
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newAccessToken, err := e.tokenGen.RegenerateAccessToken(token.RefreshToken); err == nil {
		c.JSON(http.StatusOK, gin.H{ // Return new access token as JSON with 200 status code
			"accessToken": newAccessToken,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}
}
