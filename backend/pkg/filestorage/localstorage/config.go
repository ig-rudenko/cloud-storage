package localstorage

import "os"

type Config struct {
	Path string
}

func NewConfig() *Config {
	path := os.Getenv("STORAGE_DIR")
	if path == "" {
		path = "storage"
	}
	return &Config{
		Path: path,
	}
}
