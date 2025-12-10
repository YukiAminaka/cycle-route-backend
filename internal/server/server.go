package server

import (
	"github.com/YukiAminaka/cycle-route-backend/internal/server/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}

	router.Use(cors.New(config))

	route.InitRoute(router)
	
	router.Run(":8080")
}