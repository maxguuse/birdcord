package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	DiscordToken string `env:"DISCORD_TOKEN"`
	DiscordId    string `env:"DISCORD_ID"`

	ConnectionString string `env:"CONNECTION_STRING"`
}

func New() *Config {
	cfg := &Config{}

	if env := os.Getenv("ENVIRONMENT"); env == "prod" {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			panic("Error reading config: " + err.Error())
		}
	} else {
		err := cleanenv.ReadConfig("dev.env", cfg)
		if err != nil {
			panic("Error reading config: " + err.Error())
		}
	}

	return cfg
}
