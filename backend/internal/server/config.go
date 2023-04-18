package server

import "os"

// Config ...
type Config struct {
	Address   string
	SecretKey string
}

// NewConfig ...
func NewConfig() *Config {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		key = "my-default-secret-key"
	}
	return &Config{
		Address:   ":8080",
		SecretKey: key,
	}
}
