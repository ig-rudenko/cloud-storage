package service

import (
	"fmt"
	"strings"
	"web/backend/internal/app/model"
)

// GetUserFiles Получаем список файлов и папок для пользовательского пути в хранилище
func (s *Service) GetUserFiles(user *model.User, userPath string) ([]model.FileInfo, error) {
	var err error

	// Проверяем указанный пользователем путь
	userPath, err = s.storage.ValidateUserStoragePath(user, userPath)
	if err != nil {
		return nil, err
	}

	// Получаем из хранилища данные
	userFiles, err := s.storage.ListUserFiles(userPath)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find") {
			// Если не удалось найти директорию
			return nil, fmt.Errorf("invalid path")
		}
		return nil, err
	}
	return userFiles, nil
}
