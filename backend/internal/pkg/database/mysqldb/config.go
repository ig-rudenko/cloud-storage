package mysqldb

import (
	"fmt"
	"gorm.io/gorm"
	"os"
)

type Config struct {
	DNS  string
	Gorm *gorm.Config
}

func NewConfig() *Config {

	// DSN Get database connection settings from environment variables
	var DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASS", "password"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "test_go"))

	return &Config{
		DNS:  DSN,
		Gorm: &gorm.Config{},
	}
}

// getEnv returns the value of an environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
