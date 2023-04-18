package endpoint

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web/backend/internal/model"
)

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

// RegisterNewUser 	godoc
// @Summary        	Регистрация нового пользователя.
// @Description    	Регистрация нового пользователя.
// @Tags 			auth
// @ID 				register-user
// @Accept 			json
// @Produce			json
// @Param			input body userForm true "user data"
// @Success        	201
// @Failure			400 {object} errorResponse "invalid user data"
// @Failure			500 {object} errorResponse "failed to init user storage"
// @Router         	/api/auth/register [post]
func (e *Endpoint) RegisterNewUser(c *gin.Context) {
	// Bind the JSON data to a UserForm struct
	var u userForm
	if err := c.ShouldBindJSON(&u); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Модель пользователя
	user := model.User{
		Username: u.Username,
		Password: u.Password,
	}

	// Проверяем данные пользователя
	if err := user.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Создаем пользователя, а также инициализируем новое пользовательское хранилище
	if err := e.service.InitUser(&user); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the status code 201
	c.Status(http.StatusCreated)
}

// GenerateToken 	godoc
// @Summary        	Получение JWT.
// @Description    	Создает и возвращает токен доступа и токен обновления для заданных учетных данных пользователя.
// @Tags 			auth
// @ID 				login
// @Accept 			json
// @Produce			json
// @Param			input body userForm true "user data"
// @Success        	200 {object} tokenPair "JWT Pair"
// @Failure			400 {object} errorResponse "invalid user data"
// @Failure			500 {object} errorResponse "failed to create JWT Pair"
// @Router         	/api/auth/token [post]
func (e *Endpoint) GenerateToken(c *gin.Context) {
	// Bind the JSON data to a UserForm struct
	var userFormData userForm
	if err := c.ShouldBindJSON(&userFormData); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := e.service.GetUser(userFormData.Username)
	if err != nil {
		if err.Error() == "record not found" {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if user.ComparePassword(userFormData.Password) == false {
		newErrorResponse(c, http.StatusBadRequest, "invalid password")
		return
	}

	userID := strconv.Itoa(int(user.ID))                                 // переводим ID пользователя в строку
	accessToken, refreshToken, err := e.tokenGen.CreateTokenPair(userID) // Create a token pair for user id
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("failed to create tokens: %v", err))
		return
	}

	// Return token pair as JSON with 200 status code
	c.JSON(http.StatusOK, tokenPair{accessToken, refreshToken})

}

// RefreshToken 	godoc
// @Summary        	Обновление access токена.
// @Description    	Обновляет и возвращает токен доступа для данного токена обновления.
// @Tags 			auth
// @ID 				refresh-token
// @Accept 			json
// @Produce			json
// @Param			input body refreshToken true "refresh token"
// @Success        	200 {object} accessToken "access token"
// @Failure			400,401 {object} errorResponse "invalid refresh token"
// @Router         	/api/auth/token/refresh [post]
func (e *Endpoint) RefreshToken(c *gin.Context) {
	var rt refreshToken
	if err := c.ShouldBindJSON(&rt); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if newAccessToken, err := e.tokenGen.RegenerateAccessToken(rt.Token); err == nil {
		c.JSON(http.StatusOK, accessToken{newAccessToken})
	} else {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
}
