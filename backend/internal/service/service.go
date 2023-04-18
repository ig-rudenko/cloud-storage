package service

import (
	"io"
	"io/fs"
	"mime/multipart"
	"time"
	"web/backend/internal/model"
)

type DataBase interface {
	Create(model interface{}) error
	GetOne(model interface{}, query interface{}, args ...interface{}) error
}

//type File interface {
//	fs.FileInfo
//}

type FileInfoA interface {
	Name() string
	Size() int64
	Mode() fs.FileMode
	ModTime() time.Time
	IsDir() bool
	Sys() any
}

type Storage interface {
	CreateUserStorage(name string) error
	ValidateUserStoragePath(user *model.User, path string) (validPath string, err error)
	ListUserFiles(path string) ([]fs.FileInfo, error)
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
