package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServiceName string        `yaml:"service_name" env-required:"true"`
	Env         string        `yaml:"env" env-required:"true"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	RedisPath   string        `yaml:"redis_path"`
	GRPC        GRPC          `yaml:"grpc"`
	Clients     CLientsConfig `yaml:"clients"`
}

type GRPC struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type Client struct {
	Addres       string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
}

type CLientsConfig struct {
	Event  Client `yaml:"event"`
	Ticket Client `yaml:"ticket"`
}

func New(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &Config{}, err
	}
	var config Config
	err := cleanenv.ReadConfig(filePath, &config)

	return &config, err
}
