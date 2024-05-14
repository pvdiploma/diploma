package ticketapp

import (
	"log/slog"
	"tn/internal/app"
	eventclient "tn/internal/clients/event"
	ticketgrpc "tn/internal/grpc/ticket"
	"tn/internal/services/ticket"
	"tn/internal/storage/postgresql"
	tokenmanager "tn/internal/utils/tokenManager"

	"google.golang.org/grpc"
)

type TicketApp struct {
	App *app.App
}

func NewTicketApp(log *slog.Logger, port int, storagePath string, tm *tokenmanager.TokenManager, eventClient *eventclient.Client) *TicketApp {

	storage, err := postgresql.NewStorage(storagePath)
	if err != nil {
		log.Error("Failed to create storage", err)
		panic(err)
	}
	ticketService := ticket.New(log, storage, eventClient)

	gRPCServer := grpc.NewServer()

	ticketgrpc.Register(gRPCServer, ticketService)

	return &TicketApp{
		App: app.NewApp(log, gRPCServer, port),
	}
}
