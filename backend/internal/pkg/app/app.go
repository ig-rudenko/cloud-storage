package app

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"web/backend/internal/app/endpoint"
	"web/backend/internal/app/middleware"
	"web/backend/internal/app/model"
	"web/backend/internal/app/server"
	"web/backend/internal/app/service"
	"web/backend/internal/pkg/database/mysqldb"
	"web/backend/internal/pkg/filestorage/localstorage"
	"web/backend/internal/pkg/token"
)

// App ...
type App struct {
	server *server.Server
}

// New ...
func New() (*App, error) {
	app := &App{}

	// Инициализация БД (MySQL)
	dbConfig := mysqldb.NewConfig()
	db := mysqldb.New(dbConfig)
	err := db.Open() // Подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Миграция (User)
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}

	// Создание сервера
	serverConfig := server.NewConfig()
	serverConfig.Address = ":8080"
	app.server = server.New(serverConfig)

	// Файловое хранилище (локальное)
	fileStorageConfig := localstorage.NewConfig()
	fileStorage := localstorage.New(fileStorageConfig)

	// Для создания токенов используем JWT Creator
	JWTCreator := token.New(serverConfig.SecretKey)

	// Сервис
	mainService := service.New(db, fileStorage)

	// Endpoints для обработки запросов
	endpoints := endpoint.New(mainService, JWTCreator)

	app.server.Engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	// Router для авторизации и создания нового пользователя
	authRouter := app.server.Engine.Group("/api/auth")
	authRouter.POST("/token", endpoints.GenerateToken)
	authRouter.POST("/token/refresh", endpoints.RefreshToken)
	authRouter.POST("/register", endpoints.RegisterNewUser)

	// Router для API
	apiRouter := app.server.Engine.Group("/api")
	// Middleware для проверки JWT
	apiRouter.Use(middleware.JWTAuthMiddleware(serverConfig.SecretKey))

	// Работа с файлами
	apiRouter.GET("/items/*path", endpoints.GetFilesHandler)
	apiRouter.POST("/items/upload/*path", endpoints.UploadFilesHandler)

	// Define a route for testing authorization
	apiRouter.GET("/me", func(c *gin.Context) {
		// Get user id from Gin context
		userID, _ := c.Get("user_id")

		// Return user id as JSON response
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
		})
	})

	return app, nil
}

// Run Запуск сервера
func (a *App) Run() error {
	fmt.Println("Server running")
	if err := a.server.Start(); err != nil {
		return err
	}
	return nil
}
