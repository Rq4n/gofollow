package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port   string
	Logs   LogLevel
	DB     PostgresConfig
	Mailer MailerConfig
	Auth   AuthConfig
}

type LogLevel struct {
	Style string
	Level string
}

type AuthConfig struct {
	SecretKey string
}

type MailerConfig struct {
	FromEmail string
	APIKey    string
}

type PostgresConfig struct {
	Username    string
	Password    string
	DB          string
	Host        string
	Port        string
	MaxIdleTime string
	MinIdleConn int
	MaxOpenConn int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port: os.Getenv("PORT"),
		Auth: AuthConfig{
			SecretKey: os.Getenv("JWT_SECRET"),
		},
		Logs: LogLevel{
			Style: os.Getenv("LOG_STYLE"),
			Level: os.Getenv("LOG_LEVEL"),
		},
		DB: PostgresConfig{
			Username:    os.Getenv("POSTGRES_USER"),
			Password:    os.Getenv("POSTGRES_PASSWORD"),
			Host:        os.Getenv("POSTGRES_HOST"),
			Port:        os.Getenv("POSTGRES_PORT"),
			DB:          os.Getenv("POSTGRES_DB"),
			MaxIdleTime: os.Getenv("POSTGRES_MAX_IDLE_TIME"),
			MaxOpenConn: GetInt("POSTGRES_MAX_OPEN_CONN", 0),
			MinIdleConn: GetInt("POSTGRES_MIN_IDLE_CONN", 0),
		},
		Mailer: MailerConfig{
			FromEmail: os.Getenv("MAILTRAP_FROM_EMAIL"),
			APIKey:    os.Getenv("MAILTRAP_API_KEY"),
		},
	}
	if cfg.Auth.SecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD is required")
	}
	if cfg.Mailer.APIKey == "" {
		return nil, fmt.Errorf("MAILTRAP_API_KEY is required")
	}
	return cfg, nil
}
