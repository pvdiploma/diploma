package eventapp

import (
	"log/slog"
	"tn/internal/app"
	"tn/internal/storage/postgresql"
	tokenmanager "tn/internal/utils/tokenManager"

	"google.golang.org/grpc"
)

type EventApp struct {
	App *app.App
}

func NewEventApp(log *slog.Logger, port int, storagePath string, redisPath string, tm *tokenmanager.TokenManager) *EventApp {

	storage, err := postgresql.NewStorage(storagePath)
	if err != nil {
		log.Error("Failed to create storage", err)
		panic(err)
	}
	_ = storage
	// connect to redis

	// implement event service layer

	gRPCServer := grpc.NewServer()
	// implement event grpc layer
	return &EventApp{
		App: app.NewApp(log, gRPCServer, port),
	}
}
