package config

import "os"

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
	URL         string
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
			URL:         os.Getenv("POSTGRES_URL"),
			Port:        os.Getenv("POSTGRES_PORT"),
			MaxIdleTime: os.Getenv("POSTGRES_MAX_IDLE_TIME"),
			MaxOpenConn: GetInt("POSTGRES_MAX_OPEN_CONN", 0),
			MinIdleConn: GetInt("POSTGRES_MIN_IDLE_CONN", 0),
		},
		Mailer: MailerConfig{
			FromEmail: os.Getenv("MAILTRAP_FROM_EMAIL"),
			APIKey:    os.Getenv("MAILTRAP_API_KEY"),
		},
	}
	return cfg, nil
}
