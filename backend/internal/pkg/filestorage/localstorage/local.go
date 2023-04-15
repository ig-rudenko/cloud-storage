package localstorage

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"web/backend/internal/app/model"
)

type Storage struct {
	Config *Config
}

func New(config *Config) *Storage {
	return &Storage{
		Config: config,
	}
}

func (s *Storage) CreateUserStorage(name string) error {
	if err := os.MkdirAll(s.Config.Path+"/"+name, 0755); err != nil {
		return err
	}
	return nil
}

func (s *Storage) ValidateUserStoragePath(user *model.User, path string) (validPath string, err error) {

	// Нельзя подниматься выше по директории
	if strings.Contains(path, "..") {
		return "", fmt.Errorf("invalid path")
	}

	// корневая директория пользователя
	validPath = fmt.Sprintf("%s/%d", s.Config.Path, user.ID)

	if path != "" {
		// Если path не задан, используем текущую директорию
		// Иначе, используем заданную пользователем директорию
		validPath += path
	}

	return validPath, err
}

func (s *Storage) ListUserFiles(path string) ([]model.FileInfo, error) {
	// Открываем директорию по заданному пути
	dir, err := os.Open(path)
	if err != nil {
		return []model.FileInfo{}, err
	}

	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
			return
		}
	}(dir)

	// Получаем список файлов и папок в директории
	files, err := dir.Readdir(-1)
	if err != nil {
		return []model.FileInfo{}, err
	}

	var fileInfos []model.FileInfo
	for _, file := range files {
		fileInfo := model.FileInfo{
			Name:    file.Name(),
			Size:    file.Size(),
			IsDir:   file.IsDir(),
			ModTime: file.ModTime().Format("15:04 / 01.02.2006"),
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

func (s *Storage) SaveFile(file *multipart.FileHeader, path string) error {
	log.Println(path)
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}

	dst := filepath.Join(path, file.Filename)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (s *Storage) DownloadFile(filepath string) (io.Reader, error) {
	// Открываем файл для чтения
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *Storage) CreateFolder(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

func (s *Storage) RenameElement(path, newName string) error {
	if err := os.Rename(path, filepath.Dir(path)+"/"+newName); err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteElement(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
