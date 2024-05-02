package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	eventapp "tn/internal/app/eventApp"
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

	singingKey := []byte(os.Getenv("SINGING_KEY"))

	tm := tokenmanager.NewManager(singingKey)
	eventApp := eventapp.NewEventApp(log, cfg.GRPC.Port, cfg.StoragePath, cfg.RedisPath, tm)

	go eventApp.App.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	log.Info("Stopping grpc auth server", slog.String("stop signal", sig.String()))
	eventApp.App.Stop()
	return nil

}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
