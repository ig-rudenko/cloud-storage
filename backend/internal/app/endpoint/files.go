package endpoint

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GetFilesHandler Функция для получения списка файлов и папок в указанной директории
func (e *Endpoint) GetFilesHandler(c *gin.Context) {

	userPath := c.Param("path")

	user, err := e.parseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Получаем файлы
	fileInfos, err := e.service.GetUserFiles(user, userPath)
	if err != nil {
		if strings.Contains(err.Error(), "invalid path") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Если нет файлов, то отправляем пустой массив
	if fileInfos == nil {
		c.String(http.StatusOK, "[]")
	} else {
		c.JSON(http.StatusOK, fileInfos)
	}

}
