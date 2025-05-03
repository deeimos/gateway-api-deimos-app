package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env"  env:"ENV" env-default:"local" env-required:"true"`
	HTTPConfig `yaml:"http"`
	APIs       `yaml:"client"`
}

type HTTPConfig struct {
	Port        int           `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}
type APIs struct {
	AuthAPI    API `yaml:"auth_api" env-required:"true"`
	ServersAPI API `yaml:"servers_api" env-required:"true"`
}

type API struct {
	Host    string        `yaml:"host"  env-required:"true"`
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
	UseTLS  bool          `yaml:"use_tls"  env-default:"false"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is required")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file: '%s' does not exist", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return &config
}
