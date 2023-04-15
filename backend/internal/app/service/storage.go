package service

import (
	"fmt"
	"mime/multipart"
	"strings"
	"sync"
	"web/backend/internal/app/model"
)

func (s *Service) ValidatePath(user *model.User, path string) (string, error) {
	// Проверяем указанный пользователем путь
	valid, err := s.storage.ValidateUserStoragePath(user, path)
	if err != nil {
		return "", err
	}
	return valid, nil
}

// GetUserFiles Получаем список файлов и папок для пользовательского пути в хранилище
func (s *Service) GetUserFiles(path string) ([]model.FileInfo, error) {
	var err error

	// Получаем из хранилища данные
	userFiles, err := s.storage.ListUserFiles(path)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find") {
			// Если не удалось найти директорию
			return nil, fmt.Errorf("invalid path")
		}
		return nil, err
	}
	return userFiles, nil
}

func (s *Service) UploadFiles(files []*multipart.FileHeader, path string) []error {

	// Кол-во файлов
	filesLength := len(files)

	// Создаем канал для передачи ошибок от горутин
	errCh := make(chan error, filesLength)

	// Создаем WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// Запускаем горутину для каждого файла
	for _, file := range files {
		wg.Add(1) // Увеличиваем счетчик WaitGroup

		go func(file *multipart.FileHeader) {
			defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении горутины

			if err := s.storage.SaveFile(file, path); err != nil {
				errCh <- err // Отправляем ошибку в канал, если есть
			}

		}(file)
	}

	// Создаем слайс для хранения ошибок из канала
	var errors = make([]error, 0, filesLength)

	// Запускаем горутину для сбора ошибок из канала
	go func() {
		for err := range errCh {
			errors = append(errors, err) // Добавляем ошибку в слайс, если есть
		}
	}()

	// Проверяем, есть ли ошибки в слайсе
	if len(errors) > 0 {
		return errors
	}

	// Ждем завершения всех горутин
	wg.Wait()
	close(errCh) // Закрываем канал ошибок

	return nil
}
