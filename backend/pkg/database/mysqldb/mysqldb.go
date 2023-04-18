package mysqldb

import (
	"errors"
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataBase struct {
	Config *Config
	DB     *gorm.DB
}

func New(config *Config) *DataBase {
	return &DataBase{
		Config: config,
	}
}

func (db *DataBase) Open() error {
	var err error
	db.DB, err = gorm.Open(mysql.Open(db.Config.DNS), db.Config.Gorm)
	if err != nil {
		panic("failed to connect to database")
	}
	return nil
}

func (db *DataBase) AutoMigrate(models ...interface{}) error {
	err := db.DB.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}

// handleDBEntryError Обрабатываем ошибку базы данных. Определяем текст пользователю и статус ответа
func handleDBEntryError(err error) error {
	var mysqlErr *mysqlDriver.MySQLError     // define a variable of type *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok { // check if err can be assigned to mysqlErr variable (type assertion)
		if mysqlErr.Number == 1062 { // error number for duplicate entry https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html#error_er_dup_entry
			return fmt.Errorf("a user already exists with the same username") // return from the function or do something else
		}
	}
	return err // otherwise the original error
}

func (db *DataBase) Create(model interface{}) error {
	if err := db.DB.Create(model).Error; err != nil {
		err := handleDBEntryError(err) // Обрабатываем ошибку создания нового пользователя
		return err
	}
	return nil
}

func (db *DataBase) GetOne(model interface{}, query interface{}, args ...interface{}) error {
	res := db.DB.Where(query, args...).First(model)
	if res.Error == nil {
		return nil
	}
	if res.Error == gorm.ErrRecordNotFound {
		return fmt.Errorf("record not found")
	} else {
		return res.Error
	}
}
