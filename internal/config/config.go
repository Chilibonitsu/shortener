package config

import "time"

type Config struct {
	Env                   string        `env:"ENV"`
	Port                  int           `env:"PORT"`
	DSN                   string        `env:"DSN"`
	ConfigPath            string        `env:"CONFIG_PATH"`
	DBHost                string        `env:"DB_HOST"`
	DBPort                int           `env:"DB_PORT"`
	DBUser                string        `env:"DB_USER"`
	DBPassword            string        `env:"DB_PASSWORD"`
	DBName                string        `env:"DB_NAME"`
	SSLMode               string        `env:"SSL_MODE"`
	HTTPServerAddress     string        `env:"HTTP_SERVER_ADDRESS"`
	HTTPServerTimeout     time.Duration `env:"HTTP_SERVER_TIMEOUT"`
	HTTPServerIdleTimeout time.Duration `env:"HTTP_SERVER_IDDLE_TIMEOUT"`
}
