package endpoint

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Error string `json:"error"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newStatusResponse(c *gin.Context, code int, message string) {
	c.JSON(code, statusResponse{message})
}

func newErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, errorResponse{message})
}
