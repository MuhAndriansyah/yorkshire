package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	HttpAddr         string        `env:"HTTP_ADDR" envDefault:":8080"`
	DatabaseUrl      string        `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/yorkshires?sslmode=disable"`
	RedisUrl         string        `env:"REDIS_URL" envDefault:"localhost:6379"`
	RedisPassword    string        `env:"REDIS_PASSWORD"`
	RedisDB          int           `env:"REDIS_DB" envDefault:"0"`
	LogFile          string        `env:"LOG_FILE" envDefault:"logs/app.log"`
	PrivKeyPath      string        `env:"PRIVKEY_PATH" envDefault:"files/private_key.pem"`
	PubKeyPath       string        `env:"PUBKEY_PATH" envDefault:"files/public_key.pem"`
	JWTTokenDuration time.Duration `env:"JWT_TOKEN_DURATION" envDefault:"60m"`
	JWTIssuer        string        `env:"JWT_ISSUER"`
	JWTAudience      string        `env:"JWT_AUDIENCE"`
}

func Load(filename string) (*Config, error) {
	if err := godotenv.Load(filename); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load env: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	return cfg, nil
}
