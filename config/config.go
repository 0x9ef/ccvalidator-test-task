package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      App
		Log      Log
		HTTP     HTTP
		Validate Validate
	}

	App struct {
		Name        string `env:"APP_NAME"        env-default:"ccvalidator"`
		Version     string `env:"APP_VERSION"     env-default:"0.0.1"`
		Environment string `env:"APP_ENVIRONMENT" env-default:"local"`
		BaseURL     string `env:"APP_BASE_URL"    env-default:"https://localhost:8080"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"info"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" env-default:"8080"`
	}

	Validate struct {
		AllowTestCards bool `env:"VALIDATE_ALLOW_TEST_CARDS" env-default:"true"`
	}
)

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		if err := cleanenv.ReadEnv(&config); err != nil {
			log.Fatal("failed to read env", err)
		}
	})
	return &config
}
