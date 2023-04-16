package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	app := &App{}

	// Инициализация БД (MySQL)
	dbConfig := mysqldb.NewConfig()
	db := mysqldb.New(dbConfig)
	err := db.Open() // Подключаемся к БД
	if err != nil {
		return
	}

	// Миграция (User)
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return
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

	// Документация
	app.server.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Router для авторизации и создания нового пользователя
	authRouter := app.server.Engine.Group("/api/auth")
	// router: /api/auth
	{
		authRouter.POST("/token", endpoints.GenerateToken)
		authRouter.POST("/token/refresh", endpoints.RefreshToken)
		authRouter.POST("/register", endpoints.RegisterNewUser)
	}

	// Router для API
	apiRouter := app.server.Engine.Group("/api")
	// Middleware для проверки JWT
	apiRouter.Use(middleware.JWTAuthMiddleware(serverConfig.SecretKey))
	// router: /api
	{
		// Работа с файлами
		apiRouter.GET("/items/*path", endpoints.GetFilesHandler)
		apiRouter.POST("/items/upload/*path", endpoints.UploadFilesHandler)
		apiRouter.GET("/item/*path", endpoints.DownloadFile)
		apiRouter.POST("/items/create-folder/*path", endpoints.CreateDirectory)
		apiRouter.POST("/item/rename/*path", endpoints.RenameFile)
		apiRouter.DELETE("/item/*path", endpoints.DeleteItem)
	}
	// Define a route for testing authorization
	apiRouter.GET("/me", func(c *gin.Context) {
		// Get user id from Gin context
		userID, _ := c.Get("user_id")

		// Return user id as JSON response
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
		})
	})
}
