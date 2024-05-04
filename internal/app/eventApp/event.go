package eventapp

import (
	"log/slog"
	"tn/internal/app"
	eventgrpc "tn/internal/grpc/event"
	"tn/internal/services/event"
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
	// connect to redis (later)

	eventService := event.New(log, storage, storage.DB())

	gRPCServer := grpc.NewServer()

	eventgrpc.Register(gRPCServer, eventService, tm)
	return &EventApp{
		App: app.NewApp(log, gRPCServer, port),
	}
}
