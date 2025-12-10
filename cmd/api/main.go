package main

import (
	_ "github.com/YukiAminaka/cycle-route-backend/docs"
	"github.com/YukiAminaka/cycle-route-backend/internal/server"
)

// @title           Cycle-Route API
// @version         1.0
// @description     This is a server for Cycle-Route API.
// @host            localhost:8080
// @BasePath /api/v1
func main() {
	server.Run()
}
