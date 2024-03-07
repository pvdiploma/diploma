package main

import (
	"fmt"
	"os"
	authapp "tn/internal/app/authApp"
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

	authApp := authapp.NewAuthApp(log, cfg.GRPC.Port)

	authApp.Run()
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
