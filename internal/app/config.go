package app

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type (
	Config struct {
		HTTPServer
		Cache
		Limit
	}
	HTTPServer struct {
		RunPort string `env:"SERVER_PORT" env-default:"8000"`
	}
	Cache struct {
		Host string `env:"REDIS_HOST" env-default:"localhost"`
		Port string `env:"REDIS_PORT" env-default:"6379"`
		Pass string `env:"REDIS_PASS"`
	}
	Limit struct {
		Count  int64         `env:"LIMIT" env-default:"10"`
		Period time.Duration `env:"PERIOD" env-default:"1m"`
	}
)

func NewConfig() (cfg *Config) {
	cfg = &Config{}
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	return
}
