package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"web/backend/internal/pkg/app"
)

// @title Black-Hole Cloud Storage API
// @version 1.0
// @description Simple cloud file storage.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("Starting REST API Server")
		if err := a.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Println("Shutting down server")
	a.Stop()
}
