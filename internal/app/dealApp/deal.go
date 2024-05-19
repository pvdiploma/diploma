package authapp

import (
	"log/slog"
	"tn/internal/app"
	dealgrpc "tn/internal/grpc/deal"
	"tn/internal/services/deal"
	"tn/internal/storage/postgresql"

	"google.golang.org/grpc"
)

type DealApp struct {
	App *app.App
}

func NewAuthApp(log *slog.Logger, port int, storagePath string) *DealApp {

	storage, err := postgresql.NewStorage(storagePath)
	if err != nil {
		log.Error("Failed to create storage", err)
		panic(err)
	}

	dealService := deal.New(log, storage, storage)

	gRPCServer := grpc.NewServer()

	dealgrpc.Register(gRPCServer, dealService)

	return &DealApp{
		App: app.NewApp(log, gRPCServer, port),
	}
}
