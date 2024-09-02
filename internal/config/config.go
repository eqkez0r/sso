package config

import (
	"errors"
	"flag"
	"github.com/eqkez0r/sso/internal/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

var (
	ErrEmptyConfigPath = errors.New("empty config path")
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" evn-required:"true"`
	TokenTTl    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Timeout time.Duration `yaml:"timeout"`
	Port    int           `yaml:"port"`
}

func New(l logger.Logger) *Config {
	cfg := &Config{}
	path := fetchConfigPath()
	if path == "" {
		l.Error(ErrEmptyConfigPath)
		os.Exit(1)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		l.Error(err)
		os.Exit(1)
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		l.Error(err)
		os.Exit(1)
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
