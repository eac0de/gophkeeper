package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	// Common
	ProjectName string `env:"PROJECT_NAME" envDefault:"gophkeeper"`
	IsDev       bool   `env:"IS_DEV" envDefault:"true"`

	// Server
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8081"`

	// PSQL
	PSQLHost     string `env:"PSQL_HOST" envDefault:"localhost"`
	PSQLPort     string `env:"PSQL_PORT" envDefault:"5432"`
	PSQLUsername string `env:"PSQL_USERNAME" envDefault:"postgres"`
	PSQLPassword string `env:"PSQL_PASSWORD" envDefault:"351762"`
	PSQLDBName   string `env:"PSQL_DB_NAME" envDefault:"gophkeeper"`

	// AuthService
	AuthServerAddress string `env:"AUTH_SERVER_ADDRESS" envDefault:"0.0.0.0:8080"`
}

func MustLoad() *Config {
	cfg := new(Config)
	err := env.Parse(cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return cfg
}
