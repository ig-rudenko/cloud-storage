package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"web/backend/api"
	"web/backend/config"
)

func main() {

	var err error

	// Connect to the database using the username, password, host, port and database name
	config.DB, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the user model
	err = config.DB.AutoMigrate(&api.User{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Create a Gin router with default middleware
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))
	apiRouter := router.Group("/api")
	authRouter := router.Group("/api/auth")

	// Use AuthMiddleware for routes under /api prefix
	apiRouter.Use(AuthMiddleware(config.SecretKey))

	// Define a route for generating tokens
	authRouter.POST("/token", api.GenerateToken)
	// Define a route for refreshing tokens
	authRouter.POST("/token/refresh", api.RefreshToken)
	authRouter.POST("/register", api.CreateUser)

	// Работа с файлами
	apiRouter.GET("/items/*path", api.GetFilesHandler)
	apiRouter.POST("/items/upload/*path", api.PostFilesHandler)
	apiRouter.POST("/items/create-folder/*path", api.CreateFolderHandler)

	// Работа с одним файлом
	apiRouter.GET("/item/*path", api.DownloadFile)
	apiRouter.POST("/item/rename/*path", api.RenameFileHandler)
	apiRouter.DELETE("/item/*path", api.DeleteFileHandler)

	// Define a route for testing authorization
	apiRouter.GET("/me", func(c *gin.Context) {
		// Get user id from Gin context
		userID, _ := c.Get("user_id")

		// Return user id as JSON response
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
		})
	})

	// Run router on port 8080
	log.Fatal(router.Run(":8080"))
}
