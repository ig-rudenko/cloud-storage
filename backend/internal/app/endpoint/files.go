package endpoint

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// GetFilesHandler 	godoc
// @Summary        	Список файлов.
// @Security		Bearer
// @Description    	Получение списка файлов и папок в указанной директории.
// @Tags 			storage
// @ID 				get-files-list
// @Produce			json
// @Param 			path path string true "path to directory"
// @Success        	200 {object} []model.FileInfo "user's files"
// @Failure			400 {object} errorResponse "invalid user or path"
// @Failure			500 {object} errorResponse "unable to access user storage"
// @Router         	/api/items/{path} [get]
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

// UploadFilesHandler godoc
// @Summary        	Загрузка файлов.
// @Security		Bearer
// @Description    	Сохранение одного или нескольких файлов в указанную директорию.
// @Tags 			storage
// @ID 				upload-files
// @Accept			multipart/form-data
// @Produce			json
// @Param 			path path string true "path to directory"
// @Param 			files formData file true "files"
// @Success        	200 {object} statusResponse "{n} files uploaded to {path}"
// @Failure			400 {object} errorResponse "invalid user or path"
// @Failure			500 {object} errorResponse "unable to upload"
// @Router         	/api/items/upload/{path} [post]
func (e *Endpoint) UploadFilesHandler(c *gin.Context) {

	// Получаем пользователя
	user, ok := e.parseUser(c)
	if !ok {
		return
	}

	// Проверяем указанный пользователем путь
	userPath, err := e.service.ValidatePath(user, c.Param("path"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем данные формы пользователя
	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	files := form.File["files"]

	// Сохраняем файлы файлы
	errors := e.service.UploadFiles(files, userPath)

	if len(errors) > 0 {
		var errorString string
		for _, err := range errors {
			errorString += err.Error()
		}
		newErrorResponse(c, http.StatusInternalServerError, errorString)
		return
	}

	newStatusResponse(c, http.StatusCreated, fmt.Sprintf("%d files uploaded to %s", len(files), userPath))
}

// DownloadFile 	godoc
// @Summary        	Отправка файла.
// @Security		Bearer
// @Description    	Отправляем файл пользователю.
// @Tags 			storage
// @ID 				download-files
// @Param 			path path string true "path to file"
// @Success        	200 "file"
// @Failure			400 {object} errorResponse "invalid user or path"
// @Failure			500 {object} errorResponse "unable to send user file"
// @Router         	/api/item/{path} [get]
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
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	file, err := e.service.DownloadFile(userPath)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := io.Copy(c.Writer, file); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

// CreateDirectory 	godoc
// @Summary        	Создание новой директории.
// @Security		Bearer
// @Description    	Создаем новую директорию в папке указанной в параметре path.
// @Tags 			storage
// @ID 				create-directory
// @Param 			path path string true "path to file"
// @Success        	201
// @Failure			400 {object} errorResponse "invalid user or path"
// @Failure			500 {object} errorResponse "unable to create directory"
// @Router         	/api/items/create-folder/{path} [post]
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
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := e.service.CreateFolder(userPath); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newStatusResponse(c, http.StatusCreated, fmt.Sprintf("directory %s created", userPath))
}

// RenameFile 		godoc
// @Summary        	Переименование директории/файла.
// @Security		Bearer
// @Description    	Переименновываем папку или файл указанный в параметре path.
// @Tags 			storage
// @ID 				rename-item
// @Param 			path path string true "path to dir/file"
// @Param			input body newItemName true "new file/directory name"
// @Success        	200 {object} statusResponse "renamed"
// @Failure			400 {object} errorResponse "invalid user or path"
// @Failure			500 {object} errorResponse "unable to rename directory"
// @Router         	/api/item/{path} [patch]
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

	var newName newItemName

	if err := c.ShouldBindJSON(&newName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Название нового файла должно быть указано
	if newName.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newName is required"})
		return
	}

	if err := e.service.RenameFile(newName.Name, userPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": fmt.Sprintf("%s renamed to %s", userPath, newName)})
}

// DeleteItem 		godoc
// @Summary        	Удаление директории/файла.
// @Security		Bearer
// @Description    	Удаляем папку или файл указанный в параметре path.
// @Tags 			storage
// @ID 				delete-item
// @Param 			path path string true "path to dir/file"
// @Param			input body newItemName true "new file/directory name"
// @Success        	204
// @Failure			400 {object} errorResponse "invalid user or path"
// @Failure			500 {object} errorResponse "unable to delete directory"
// @Router         	/api/item/{path} [delete]
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
	c.Status(http.StatusNoContent)
}
