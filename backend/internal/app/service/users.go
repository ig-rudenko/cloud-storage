package service

import (
	"strconv"
	"web/backend/internal/app/model"
)

func (s *Service) InitUser(user *model.User) error {

	// Шифруем пароль пользователя
	if err := user.EncryptPassword(); err != nil {
		return err
	}

	// Сохраняем пользователя в БД
	if err := s.db.Create(user); err != nil {
		return err
	}

	userID := strconv.FormatUint(uint64(user.ID), 10) // переводим ID пользователя в строку

	// Создаем хранилище для нового пользователя
	if err := s.storage.CreateUserStorage(userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUser(username string) (model.User, error) {
	var user model.User
	if err := s.db.GetOne(&user, "username = ?", username); err != nil {
		return model.User{}, err
	}
	return user, nil
}
