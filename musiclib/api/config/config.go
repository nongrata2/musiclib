package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServerAddress string        `yaml:"http_server_address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8081"`
	HttpServerTimeout time.Duration `yaml:"http_server_timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
	LogLevel          string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	DBHost    string `yaml:"db_host" env:"DB_HOST" env-default:"db"`
	DBUser    string `yaml:"db_user" env:"DB_USER" env-default:"postgres"`
	DBPassword    string `yaml:"db_password" env:"DB_PASSWORD" env-default:"postgres"`
	DBName    string `yaml:"db_name" env:"DB_NAME" env-default:"postgres"`
	DBPort    string `yaml:"db_port" env:"DB_PORT" env-default:"5432"`
}


func MustLoadCfg(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read http config %q: %s", configPath, err)
	}
	return cfg
}
