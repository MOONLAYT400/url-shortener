package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	Clients ClientsConfig `yaml:"clients"`
	// AppSecret string `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
}

type 	HTTPServer  struct {
		Address     string  `yaml:"address" env-default:"localhost:8000"`
		Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
		IdleTimeout time.Duration     `yaml:"idle_timeout" env-default:"60s"`
		User string `yaml:"user" env-required:"true"`
		Password string `yaml:"password" env-required:"true" env:"USER_PASSWORD"`
	}

	type CLient struct{
		Address string `yaml:"address"`
		Timeout time.Duration `yaml:"timeout"`
		RetriesCount int `yaml:"retries_count"`
	}

	type ClientsConfig struct{
		SSO CLient `yaml:"sso"`
	}

	func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatal("CONFIG_PATH is not set")
		}

		//check if file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file does not exist: %s", configPath)
		}

		var cfg Config
		if err:= cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("error reading config file: %s", err)
		}

		return &cfg
	}