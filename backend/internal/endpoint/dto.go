package endpoint

// tokenPair пара токенов
type tokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// accessToken
type accessToken struct {
	Token string `json:"accessToken"`
}

// refreshToken
type refreshToken struct {
	Token string `json:"refreshToken"`
}

// newName новое название для файла или директории
type newItemName struct {
	Name string `json:"newName"`
}

// userForm форма для данных пользователя
type userForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// File is a struct that holds information about a file or folder
type userFile struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime string `json:"modTime"`
}
