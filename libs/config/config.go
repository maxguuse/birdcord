package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	EnvProduction = "prod"
)

type Config struct {
	Environment string `env:"ENVIRONMENT"`
	Version     string `env:"VERSION"`

	DiscordToken string `env:"DISCORD_TOKEN"`
	DiscordId    string `env:"DISCORD_ID"`

	ConnectionString string `env:"CONNECTION_STRING"`

	RedisDSN      string `env:"REDIS_DSN"`
	RedisPassword string `env:"REDIS_PASS"`

	DebugGuildId string `env:"DEBUG_GUILD_ID"`
}

func New() *Config {
	cfg := &Config{}

	if env := os.Getenv("ENVIRONMENT"); env == EnvProduction {
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
