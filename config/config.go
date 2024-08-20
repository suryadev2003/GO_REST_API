package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	} `yaml:"database"`
	JWT struct {
		Secret     string `yaml:"secret"`
		Expiration string `yaml:"expiration"`
	} `yaml:"jwt"`
}

func LoadConfig() *Config {
	var config Config

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, falling back to config.yaml")
	}

	if !viper.IsSet("DB_USER") {
		file, err := os.Open("config.yaml")
		if err != nil {
			log.Fatalf("Error opening config.yaml file: %v", err)
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			log.Fatalf("Error decoding config.yaml file: %v", err)
		}

		viper.SetDefault("DB_USER", config.Database.User)
		viper.SetDefault("DB_PASSWORD", config.Database.Password)
		viper.SetDefault("DB_NAME", config.Database.Name)
		viper.SetDefault("DB_HOST", config.Database.Host)
		viper.SetDefault("DB_PORT", config.Database.Port)
		viper.SetDefault("JWT_SECRET", config.JWT.Secret)
		viper.SetDefault("JWT_EXPIRATION", config.JWT.Expiration)
		viper.SetDefault("SERVER_PORT", config.Server.Port)
	}

	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")
	config.Database.Name = viper.GetString("DB_NAME")
	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetInt("DB_PORT")
	config.JWT.Secret = viper.GetString("JWT_SECRET")
	config.JWT.Expiration = viper.GetString("JWT_EXPIRATION")
	config.Server.Port = viper.GetInt("SERVER_PORT")

	return &config
}
