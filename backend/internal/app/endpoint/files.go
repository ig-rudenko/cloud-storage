package endpoint

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

	if len(errors) > 0 {
		for _, err := range errors {
			c.Error(err)
		}
		c.JSON(http.StatusInternalServerError, c.Errors)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": fmt.Sprintf("%d files uploaded to %s", len(files), userPath)})
}

// DownloadFile Отправляем файл пользователю
func (e *Endpoint) DownloadFile(c *gin.Context) {
	userPath := c.Param("path")

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

	file, err := e.service.DownloadFile(userPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := io.Copy(c.Writer, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// CreateDirectory Создаем новую директорию
func (e *Endpoint) CreateDirectory(c *gin.Context) {
	userPath := c.Param("path")

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

	if err := e.service.CreateFolder(userPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": fmt.Sprintf("directory %s created", userPath)})
}

// RenameFile Создаем новую директорию
func (e *Endpoint) RenameFile(c *gin.Context) {
	userPath := c.Param("path")

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

	var data map[string]string

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newName := data["newName"]

	// Название нового файла должно быть указано
	if newName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newName is required"})
		return
	}

	if err := e.service.RenameFile(newName, userPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": fmt.Sprintf("%s renamed to %s", userPath, newName)})
}

func (e *Endpoint) DeleteItem(c *gin.Context) {
	userPath := c.Param("path")

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

	if err := e.service.DeleteFile(userPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": fmt.Sprintf("%s deleted", userPath)})
}
