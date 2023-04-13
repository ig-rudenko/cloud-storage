package service

import "web/backend/internal/app/model"

type DataBase interface {
	Create(model interface{}) error
	GetOne(model interface{}, query interface{}, args ...interface{}) error
}

type Storage interface {
	CreateUserStorage(name string) error
	ValidateUserStoragePath(user *model.User, path string) (validPath string, err error)
	ListUserFiles(path string) ([]model.FileInfo, error)
}

type Service struct {
	db      DataBase
	storage Storage
}

func New(db DataBase, storage Storage) *Service {
	return &Service{
		db:      db,
		storage: storage,
	}
}