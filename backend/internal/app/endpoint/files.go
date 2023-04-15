package endpoint

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

// GetFilesHandler Функция для получения списка файлов и папок в указанной директории
func (e *Endpoint) GetFilesHandler(c *gin.Context) {

	user, ok := e.parseUser(c)
	if !ok {
		return
	}

	// Проверяем указанный пользователем путь
	userPath, err := e.service.ValidatePath(user, c.Param("path"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Получаем файлы
	fileInfos, err := e.service.GetUserFiles(userPath)
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

// UploadFilesHandler Функция для загрузки одного или нескольких файлов в указанную директорию
func (e *Endpoint) UploadFilesHandler(c *gin.Context) {

	// Получаем пользователя
	user, ok := e.parseUser(c)
	if !ok {
		return
	}

	// Проверяем указанный пользователем путь
	userPath, err := e.service.ValidatePath(user, c.Param("path"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем данные формы пользователя
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]

	// Сохраняем файлы файлы
	errors := e.service.UploadFiles(files, userPath)

	if len(errors) > 1 {
		// Если ошибок больше одной, то отправляем слайс текстов ошибок
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": func() []string {
					var e []string
					for _, err := range errors {
						e = append(e, err.Error())
					}
					return e
				}(),
			},
		)
		return
	} else if len(errors) == 1 {
		// Если ошибка одна, то отправляем её текст
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors[0].Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": fmt.Sprintf("%d files uploaded to %s", len(files), userPath)})
}

// DownloadFile Отправляем файл пользователю
func (e *Endpoint) DownloadFile(c *gin.Context) {
	userPath := c.Param("path")

	//user, err := e.parseUser(c)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	//	return
	//}
	//
	// Проверьте правильность пути и существование файла
	fileName := filepath.Base(userPath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(fileName)
}
