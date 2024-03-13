package authapp

import (
	"log/slog"
	"tn/internal/app"
	authgrpc "tn/internal/grpc/auth"

	"google.golang.org/grpc"
)

func NewAuthApp(log *slog.Logger, auth authgrpc.Auth, port int) *app.App {

	// implement interface service

	//start grpc server
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, auth)

	return app.NewApp(log, gRPCServer, port)
}
