package config

import (
	"fmt"
	"gorm.io/gorm"
	"os"
)

// getEnv returns the value of an environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// StorageDir Директория для хранения файлов пользователей
var StorageDir = getEnv("STORAGE_DIR", "storage")

// DB is a global variable that holds the database connection
var DB *gorm.DB

// SecretKey is a global variable that holds the secret key for signing and validating tokens
var SecretKey = getEnv("SECRET_KEY", "mysecretkey")

// DSN Get database connection settings from environment variables
var DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	getEnv("DB_USER", "root"),
	getEnv("DB_PASS", "password"),
	getEnv("DB_HOST", "localhost"),
	getEnv("DB_PORT", "3306"),
	getEnv("DB_NAME", "test_go"))
