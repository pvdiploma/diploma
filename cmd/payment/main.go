package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	paymentapp "tn/internal/app/paymentApp"
	eventclient "tn/internal/clients/event"
	ticketclient "tn/internal/clients/ticket"

	dealclient "tn/internal/clients/deal"
	"tn/internal/config"
	tokenmanager "tn/internal/utils/tokenManager"
	"tn/pkg/logger"

	"github.com/joho/godotenv"
)

func run() error {
	err := godotenv.Load("local.env")
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
	}

	configPathAuth := os.Getenv("TN_CONFIG_PATH_EVENT")

	cfg, err := config.New(configPathAuth)
	if err != nil {
		return err
	}

	log, err := logger.SetupLogger(cfg)
	if err != nil {
		return err
	}
	// why i dont use ticket client???
	singingKey := []byte(os.Getenv("SINGING_KEY"))

	tm := tokenmanager.NewManager(singingKey)

	eventClient, err := eventclient.NewClient(
		context.Background(),
		log,
		cfg.Clients.Event.Addres,
		cfg.Clients.Event.Timeout,
		cfg.Clients.Event.RetriesCount,
	)
	if err != nil {
		log.Error("Failed to create event client", logger.Err(err))
		return err
	}

	ticketClient, err := ticketclient.NewClient(
		context.Background(),
		log,
		cfg.Clients.Ticket.Addres,
		cfg.Clients.Ticket.Timeout,
		cfg.Clients.Ticket.RetriesCount,
	)
	if err != nil {
		log.Error("Failed to create ticket client", logger.Err(err))
		return err
	}

	dealClient, err := dealclient.NewClient(
		context.Background(),
		log,
		cfg.Clients.Deal.Addres,
		cfg.Clients.Deal.Timeout,
		cfg.Clients.Deal.RetriesCount,
	)

	paymentApp := paymentapp.NewEventApp(log, cfg.GRPC.Port, cfg.StoragePath, tm, eventClient, ticketClient, dealClient)

	go paymentApp.App.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	log.Info("Stopping grpc auth server", slog.String("stop signal", sig.String()))
	paymentApp.App.Stop()
	return nil

}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
