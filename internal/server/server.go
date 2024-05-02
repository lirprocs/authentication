package server

import (
	"aut_reg/database"
	"context"
	"errors"
	"fmt"
	pb "github.com/lirprocs/protosSSO/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.AuthServer
	authService Auth
}

func NewServer() *Server {
	s := &Server{}
	return s
}

type Auth interface {
	RegisterUser(ctx context.Context, email, username, password string) (int64, error)
	LoginUser(ctx context.Context, username, password string) (string, error)
}

func Register(grpcServer *grpc.Server, authService Auth) {
	pb.RegisterAuthServer(grpcServer, &Server{authService: authService})
}

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	if in.Email == "" {
		//TODO
		return nil, fmt.Errorf("Pleas enter Email")
	}
	if in.Username == "" {
		//TODO
		return nil, fmt.Errorf("Pleas enter Username")
	}
	if in.Password == "" {
		//TODO
		return nil, fmt.Errorf("Pleas enter Password")
	}

	id, err := s.authService.RegisterUser(ctx, in.GetEmail(), in.GetUsername(), in.GetPassword())

	if err != nil {
		if errors.Is(err, database.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &pb.RegisterResponse{UserId: id}, nil
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if in.Username == "" {
		return nil, fmt.Errorf("Pleas enter Username")
	}
	if in.Password == "" {
		return nil, fmt.Errorf("Pleas enter Password")
	}

	token, err := s.authService.LoginUser(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid username or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &pb.LoginResponse{Token: token}, nil
}
