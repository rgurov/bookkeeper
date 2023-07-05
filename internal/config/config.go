package config

import (
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Host             string        `env:"HOST" env-default:"127.0.0.1"`
	Port             string        `env:"PORT" env-default:"4567"`
	Timeout          time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout      time.Duration `env:"IDLE_TIMEOUT" env-default:"30s"`
	JwtSecret        string        `env:"JWT_SECRET" env-default:"secret"`
	PostgresHost     string        `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort     string        `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string        `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword string        `env:"POSTGRES_PASSWORD" env-default:"123"`
	PostgresDatabase string        `env:"POSTGRES_DB" env-default:"bookkeeper"`
}

var readOnce = sync.Once{}
var cfg Config

func Read() *Config {
	readOnce.Do(func() {
		_ = godotenv.Load()
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			panic(err)
		}
	})
	return &cfg
}
