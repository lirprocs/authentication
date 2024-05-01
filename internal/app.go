package app

import (
	"aut_reg/database"
	grpcapp "aut_reg/internal/grpc"
	"aut_reg/internal/services"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(port int, storagePath string, tokenTTL time.Duration) *App {
	db, err := database.InitDB(storagePath)
	if err != nil {
		// TODO
		panic(err)
	}
	storage := services.New(db, db, tokenTTL)
	server := grpcapp.New(port, storage)

	return &App{GRPCServer: server}
}
