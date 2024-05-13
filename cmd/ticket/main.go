package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	ticketapp "tn/internal/app/ticketApp"
	eventclient "tn/internal/clients/event"
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

	configPathAuth := os.Getenv("TN_CONFIG_PATH_SERVICE")

	cfg, err := config.New(configPathAuth)
	if err != nil {
		return err
	}

	log, err := logger.SetupLogger(cfg)
	if err != nil {
		return err
	}

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

	ticketApp := ticketapp.NewTicketApp(log, cfg.GRPC.Port, cfg.StoragePath, tm, eventClient)

	go ticketApp.App.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	log.Info("Stopping grpc auth server", slog.String("stop signal", sig.String()))
	ticketApp.App.Stop()
	return nil

}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
