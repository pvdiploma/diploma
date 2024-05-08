package ticketapp

import (
	"log/slog"
	"tn/internal/app"
	"tn/internal/storage/postgresql"
	tokenmanager "tn/internal/utils/tokenManager"
)

type TicketApp struct {
	App *app.App
}

func NewTicketApp(log *slog.Logger, port int, storagePath string, tm *tokenmanager.TokenManager) *TicketApp {

	// init storage
	storage, err := postgresql.NewStorage(storagePath)
	if err != nil {
		log.Error("Failed to create storage", err)
		panic(err)
	}
	//ticketService
	// init grpc server
	//gRPCServer := grpc.NewServer()
	_ = storage

	return &TicketApp{
		App: app.NewApp(log, nil, port),
	}
}
