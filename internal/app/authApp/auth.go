package authapp

import (
	"log/slog"
	"tn/internal/app"
	"tn/internal/grpc/auth"

	"google.golang.org/grpc"
)

func NewAuthApp(log *slog.Logger, port int) *app.App {

	// implement interface service

	//start grpc server
	gRPCServer := grpc.NewServer()
	auth.Register(gRPCServer)

	return app.NewApp(log, gRPCServer, port)
}
