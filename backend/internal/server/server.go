package server

import "github.com/gin-gonic/gin"

// Server ...
type Server struct {
	Config *Config
	Engine *gin.Engine
}

// New ...
func New(config *Config) *Server {
	return &Server{
		Config: config,
		Engine: gin.Default(),
	}
}

// Start ...
func (s Server) Start() error {
	err := s.Engine.Run(s.Config.Address)
	if err != nil {
		return err
	}
	return nil
}
