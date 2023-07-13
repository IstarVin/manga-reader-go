package server

import (
	"github.com/gin-gonic/gin"
)

// NewServer creates a new server for models
func NewServer() *gin.Engine {
	// Setup gin engine
	//gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	api := engine.Group("/api/v1")
	{
		mangaAPI(api)
		categoryAPI(api)
	}

	return engine
}
