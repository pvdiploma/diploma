package main

import (
	"fmt"
	"os"
	"sso/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
	}

	configPathAuth := os.Getenv("TN_CONFIG_PATH_AUTH")
	cfg, err := config.New(configPathAuth)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)
}
