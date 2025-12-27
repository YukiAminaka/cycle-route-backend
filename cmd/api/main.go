package main

import (
	"context"
	"os"

	"github.com/YukiAminaka/cycle-route-backend/config"
	_ "github.com/YukiAminaka/cycle-route-backend/docs"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/server"
)

// @title			Cycle-Route API
// @version		1.0
// @description	This is a server for Cycle-Route API.
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := config.GetConfig()
	pool := database.NewDB(conf.DB)
	defer pool.Close()

	q := dbgen.New(pool)
	if err := server.Run(ctx, conf, q); err != nil {
		os.Exit(1)
	}
}
