package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary  Primary        `koanf:"primary" validate:"required"`
	Server   ServerConfig   `koanf:"server" validate:"required"`
	Auth     AuthConfig     `koanf:"auth" validate:"required"`
	Database DatabaseConfig `koanf:"database" validate:"required"`
}

type Primary struct {
	Env      string `koanf:"env" validate:"required"`
	LogLevel string `koanf:"loglevel" validate:"required"`
}

type ServerConfig struct {
	Port string `koanf:"port" validate:"required"`
}

type AuthConfig struct {
	SecretKey string `koanf:"secretkey" validate:"required"`
}

type DatabaseConfig struct {
	Host     string `koanf:"host" validate:"required"`
	Port     int    `koanf:"port" validate:"required"`
	User     string `koanf:"user" validate:"required"`
	Password string `koanf:"password" validate:"required"`
	Name     string `koanf:"name" validate:"required"`
	SSLMode  string `koanf:"sslmode" validate:"required"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	k := koanf.New(".")
	// envToConfigKey transforms environment variables (e.g. IKIRU_SERVER_PORT)
	// into koanf keys (e.g. server.port).
	envToConfigKey := func(s string) string {
		s = strings.TrimPrefix(s, AppName+"_")
		s = strings.ToLower(s)
		return strings.ReplaceAll(s, "_", ".")
	}

	if err := k.Load(env.Provider(AppName+"_", ".", envToConfigKey), nil); err != nil {
		logger.Fatal().Err(err).Msg("failed to load environment variables")
	}
	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		logger.Fatal().Err(err).Msg("failed to unmarshal environment variables")
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		logger.Fatal().Err(err).Msg("failed to validate environment variables")
	}

	return &config, nil
}
