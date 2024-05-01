package grpcapp

import (
	"aut_reg/internal/server"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type App struct {
	GRPCServer *grpc.Server
	port       int
}

func New(port int, authService server.Auth) *App {
	gRPCServer := grpc.NewServer()
	server.Register(gRPCServer, authService)

	return &App{
		GRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", a.port))
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}
	if err = a.GRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}

func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		//TODO
		log.Fatalf("Error run server")
	}
}

func (a *App) Stop() {
	a.GRPCServer.GracefulStop()
}
