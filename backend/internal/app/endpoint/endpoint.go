package endpoint

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"web/backend/internal/app/model"
)

type Service interface {
	InitUser(user *model.User) error
	GetUser(username string) (model.User, error)
	GetUserFiles(user *model.User, userPath string) ([]model.FileInfo, error)
}

type TokenCreator interface {
	RegenerateAccessToken(string) (string, error)
	CreateTokenPair(string) (string, string, error)
	CreateAccessToken(string) (string, error)
	CreateRefreshToken(string) (string, error)
}

type Endpoint struct {
	service  Service
	tokenGen TokenCreator
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func New(svr Service, token TokenCreator) *Endpoint {
	return &Endpoint{
		service:  svr,
		tokenGen: token,
	}
}

func (e *Endpoint) parseUser(c *gin.Context) (*model.User, error) {
	// Get user id from Gin context
	userID, _ := c.Get("user_id")
	if userID == nil {
		return nil, fmt.Errorf("invalid user")
	}

	// Конвертируем ID пользователя из строки в int
	id, err := strconv.Atoi(userID.(string))
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	// Создаем пользователя
	user := &model.User{
		ID: uint(id),
	}
	return user, nil
}
