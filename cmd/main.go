package main

import (
	"aut_reg/internal"
	"aut_reg/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	server := app.New(cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL) //44044 | "./database/users.db"

	go server.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	server.GRPCServer.Stop()
}

// go run cmd/main.go -config ./config/test.yaml
