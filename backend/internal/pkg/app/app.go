package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "web/backend/docs"
	"web/backend/internal/endpoint"
	"web/backend/internal/middleware"
	"web/backend/internal/model"
	"web/backend/internal/server"
	"web/backend/internal/service"
	"web/backend/pkg/database/mysqldb"
	"web/backend/pkg/filestorage/localstorage"
	"web/backend/pkg/token/jwt"
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
	JWTCreator := jwt.New(serverConfig.SecretKey)

	// Сервис
	mainService := service.New(db, fileStorage)

	// Endpoints для обработки запросов
	endpoints := endpoint.New(mainService, JWTCreator)

	//app.server.Engine.Use(cors.New(cors.Config{
	//	AllowOrigins: []string{"*"},
	//	AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
	//	AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	//}))

	// Документация
	app.server.Engine.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		apiRouter.POST("/items/create-folder/*path", endpoints.CreateDirectory)
		apiRouter.GET("/item/*path", endpoints.DownloadFile)
		apiRouter.PATCH("/item/*path", endpoints.RenameFile)
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

func (a *App) Stop() {
	// Остановка
}
