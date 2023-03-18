package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileInfo is a struct that holds information about a file or folder
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime string `json:"modTime"`
}

func userPathHandler(c *gin.Context) string {
	// Получаем параметр path из URL
	path := c.Param("path")

	// Нельзя подниматься выше по директории
	if strings.Contains(path, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return ""
	}

	// Get user id from Gin context
	userID, _ := c.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return ""
	}
	if path == "" {
		// Если path не задан, используем текущую директорию
		path = fmt.Sprintf("storage/%s", userID)
	} else {
		// Иначе, используем заданную пользователем директорию
		path = fmt.Sprintf("storage/%s%s", userID, path)
	}
	return path
}

// Функция для получения списка файлов и папок в указанной директории
func getFilesHandler(c *gin.Context) {
	path := userPathHandler(c)
	if path == "" {
		return
	}

	// Открываем директорию по заданному пути
	dir, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}(dir)

	// Получаем список файлов и папок в директории
	files, err := dir.Readdir(-1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fileInfos []FileInfo
	for _, file := range files {
		fileInfo := FileInfo{
			Name:    file.Name(),
			Size:    file.Size(),
			IsDir:   file.IsDir(),
			ModTime: file.ModTime().Format("15:04:05 01-02-2006"),
		}
		fileInfos = append(fileInfos, fileInfo)
	}
	c.JSON(http.StatusOK, fileInfos)
}

func downloadFile(c *gin.Context) {
	path := userPathHandler(c)
	if path == "" {
		return
	}
	// Проверьте правильность пути и существование файла
	fileName := filepath.Base(path)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(path)
}

// Функция для загрузки одного или нескольких файлов в указанную директорию
func postFilesHandler(c *gin.Context) {
	path := userPathHandler(c)
	if path == "" {
		return
	}

	// Проверьте правильность пути и создайте директорию при необходимости
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	for _, file := range files {
		// Сохраните файл в директории path с оригинальным именем
		dst := filepath.Join(path, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"status": fmt.Sprintf("%d files uploaded to %s", len(files), path)})
}

func createFolderHandler(c *gin.Context) {
	path := userPathHandler(c)
	if path == "" {
		return
	}

	// Проверьте правильность пути
	err := os.MkdirAll(path, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": fmt.Sprintf("Directory %s created", path)})
}

// Функция для перемещения или переименования одного или нескольких файлов или папок
func renameFileHandler(c *gin.Context) {
	path := userPathHandler(c)
	if path == "" {
		return
	}

	// Проверьте правильность пути и существование файла или директории
	newName := c.PostForm("newName")
	if newName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing newName parameter"})
		return
	}

	err := os.Rename(path, filepath.Dir(path)+"/"+newName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("%s renamed to %s", path, newName)})
}

// Функция для удаления одного или нескольких файлов или папок
func deleteFileHandler(c *gin.Context) {
	path := userPathHandler(c)
	if path == "" {
		return
	}

	if err := os.Remove(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"status": fmt.Sprintf("%s deleted", path)})
}
