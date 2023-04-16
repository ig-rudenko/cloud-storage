package endpoint

import (
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"web/backend/internal/app/model"
)

type Service interface {
	InitUser(user *model.User) error
	GetUser(username string) (model.User, error)
	ValidatePath(user *model.User, path string) (string, error)
	GetUserFiles(userPath string) ([]model.FileInfo, error)
	UploadFiles(files []*multipart.FileHeader, path string) []error
	DownloadFile(string) (io.Reader, error)
	CreateFolder(filepath string) error
	RenameFile(filepath, newName string) error
	DeleteFile(filepath string) error
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

// tokenPair пара токенов
type tokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// accessToken
type accessToken struct {
	Token string `json:"accessToken"`
}

// refreshToken
type refreshToken struct {
	Token string `json:"refreshToken"`
}

// newName новое название для файла или директории
type newItemName struct {
	Name string `json:"newName"`
}

// userForm форма для данных пользователя
type userForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(svr Service, token TokenCreator) *Endpoint {
	return &Endpoint{
		service:  svr,
		tokenGen: token,
	}
}

func (e *Endpoint) parseUser(c *gin.Context) (*model.User, bool) {
	// Get user id from Gin context
	userID, _ := c.Get("user_id")
	if userID == nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user")
		return nil, false
	}

	// Конвертируем ID пользователя из строки в int
	id, err := strconv.Atoi(userID.(string))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return nil, false
	}

	// Создаем пользователя
	user := &model.User{
		ID: uint(id),
	}
	return user, true
}
