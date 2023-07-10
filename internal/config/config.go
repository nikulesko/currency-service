package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8088"`
	TimeOut time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"http_server"`
}

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"idle_timeout"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("There is no CONFIG_PATH environment variable")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}