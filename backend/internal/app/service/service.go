package service

import (
	"io"
	"mime/multipart"
	"web/backend/internal/app/model"
)

type DataBase interface {
	Create(model interface{}) error
	GetOne(model interface{}, query interface{}, args ...interface{}) error
}

type Storage interface {
	CreateUserStorage(name string) error
	ValidateUserStoragePath(user *model.User, path string) (validPath string, err error)
	ListUserFiles(path string) ([]model.FileInfo, error)
	SaveFile(file *multipart.FileHeader, path string) error
	DownloadFile(string) (io.Reader, error)
	DeleteElement(path string) error
	RenameElement(path, newName string) error
	CreateFolder(path string) error
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
