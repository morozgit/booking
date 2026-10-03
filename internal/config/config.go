package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env              string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	HTTPServerConfig `yaml:"http_server"`
	DatabaseConfig   `yaml:"database"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     uint16 `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
}

func MustLoad() *Config {
	_ = godotenv.Load()
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	return &cfg
}
