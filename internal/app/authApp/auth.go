package authapp

import (
	"log/slog"
	"tn/internal/app"
	authgrpc "tn/internal/grpc/auth"
	"tn/internal/services/auth"
	"tn/internal/storage/postgresql"
	tokenmanager "tn/internal/utils/tokenManager"

	"google.golang.org/grpc"
)

type AuthApp struct {
	App *app.App
}

func NewAuthApp(log *slog.Logger, port int, storagePath string, tm *tokenmanager.TokenManager) *AuthApp {

	storage, err := postgresql.NewStorage(storagePath)
	if err != nil {
		log.Error("Failed to create storage", err)
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tm)

	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &AuthApp{
		App: app.NewApp(log, gRPCServer, port),
	}
}
