package paymentapp

import (
	"log/slog"
	"tn/internal/app"
	"tn/internal/storage/postgresql"
	tokenmanager "tn/internal/utils/tokenManager"

	"google.golang.org/grpc"
)

type PaymentApp struct {
	App *app.App
}

func NewEventApp(log *slog.Logger, port int, storagePath string, tm *tokenmanager.TokenManager) *PaymentApp {

	storage, err := postgresql.NewStorage(storagePath)
	if err != nil {
		log.Error("Failed to create storage", err)
		panic(err)
	}
	_ = storage
	// eventService := event.New(log, storage, storage.DB())

	gRPCServer := grpc.NewServer()

	// eventgrpc.Register(gRPCServer, eventService, tm)
	return &PaymentApp{
		App: app.NewApp(log, gRPCServer, port),
	}
}
