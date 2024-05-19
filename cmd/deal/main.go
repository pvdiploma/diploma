package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	dealapp "tn/internal/app/dealApp"
	"tn/internal/config"
	"tn/pkg/logger"

	"github.com/joho/godotenv"
)

func run() error {
	err := godotenv.Load("local.env")
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
	}

	configPathAuth := os.Getenv("TN_CONFIG_PATH_AUTH")

	cfg, err := config.New(configPathAuth)
	if err != nil {
		return err
	}

	log, err := logger.SetupLogger(cfg)
	if err != nil {
		return err
	}

	dealApp := dealapp.NewAuthApp(log, cfg.GRPC.Port, cfg.StoragePath)

	go dealApp.App.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	log.Info("Stopping grpc auth server", slog.String("stop signal", sig.String()))
	dealApp.App.Stop()
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
