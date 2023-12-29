package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"local"`
	StorageUrl string        `yaml:"storage_url" env-required:"true"`
	TokenTTL   time.Duration `yaml:"token_ttl"`
	GRPC       GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("Config path empty or env variable doesnt set!")
	}

	config, err := parseConfig(configPath)
	if err != nil {
		panic("Cant parse config: " + err.Error())
	}

	return config
}

func fetchConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")

	return configPath
}

func parseConfig(configPath string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	return &cfg, err
}
