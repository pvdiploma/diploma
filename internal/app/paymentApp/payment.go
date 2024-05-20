package paymentapp

import (
	"log/slog"
	"tn/internal/app"
	eventclient "tn/internal/clients/event"
	ticketclient "tn/internal/clients/ticket"
	paymentgrpc "tn/internal/grpc/payment"
	"tn/internal/services/payment"
	"tn/internal/storage/postgresql"
	tokenmanager "tn/internal/utils/tokenManager"

	"google.golang.org/grpc"
)

type PaymentApp struct {
	App *app.App
}

func NewEventApp(
	log *slog.Logger,
	port int,
	storagePath string,
	tm *tokenmanager.TokenManager,
	eventClient *eventclient.Client,
	ticketClient *ticketclient.Client,
) *PaymentApp {

	storage, err := postgresql.NewStorage(storagePath)
	if err != nil {
		log.Error("Failed to create storage", err)
		panic(err)
	}
	_ = storage // for future
	paymentService := payment.New(log, ticketClient, eventClient)

	gRPCServer := grpc.NewServer()

	paymentgrpc.Register(gRPCServer, paymentService, tm)
	return &PaymentApp{
		App: app.NewApp(log, gRPCServer, port),
	}
}
