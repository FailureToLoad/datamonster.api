package main

import (
	"context"

	"github.com/failuretoload/datamonster/server"
	"github.com/failuretoload/datamonster/settlement"
	postgres "github.com/failuretoload/datamonster/store/postgres"
	"github.com/failuretoload/datamonster/survivor"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	connPool   *pgxpool.Pool
	app        server.Server
	appContext context.Context
)

func init() {
	appContext = context.Background()
	connPool = postgres.InitConnPool(appContext)
	app = server.NewServer()
}

func main() {
	defer connPool.Close()
	settlementController := settlement.NewController(connPool)
	survivorController := survivor.NewController(connPool)
	survivorController.RegisterRoutes(app.Mux)
	settlementController.RegisterRoutes(app.Mux)
	app.Run()
}
