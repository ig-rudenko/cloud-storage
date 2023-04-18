package endpoint

import (
	"io"
	"io/fs"
	"mime/multipart"
	"web/backend/internal/model"
)

type Service interface {
	InitUser(user *model.User) error
	GetUser(username string) (model.User, error)
	ValidatePath(user *model.User, path string) (string, error)
	GetUserFiles(userPath string) ([]fs.FileInfo, error)
	UploadFiles(files []*multipart.FileHeader, path string) []error
	DownloadFile(string) (io.Reader, error)
	CreateFolder(filepath string) error
	RenameFile(filepath, newName string) error
	DeleteFile(filepath string) error
}

type TokenCreator interface {
	RegenerateAccessToken(string) (string, error)
	CreateTokenPair(string) (string, string, error)
	CreateAccessToken(string) (string, error)
	CreateRefreshToken(string) (string, error)
}

type Endpoint struct {
	service  Service
	tokenGen TokenCreator
}

func New(svr Service, token TokenCreator) *Endpoint {
	return &Endpoint{
		service:  svr,
		tokenGen: token,
	}
}
