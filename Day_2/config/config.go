package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBUser   string `mapstructure:"DB_USER"`
	DBPass   string `mapstructure:"DB_PASS"`
	DBName   string `mapstructure:"DB_NAME"`
	DBHost   string `mapstructure:"DB_HOST"`
	DBPort   int    `mapstructure:"DB_PORT"`
	APPPort  int    `mapstructure:"APP_PORT"`
}

func load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./.development")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("error reading config file ", err)
		return &Config{}
	}
	return &Config{
		DBDriver: viper.GetString("DB_DRIVER"),
		DBUser:   viper.GetString("DB_USER"),
		DBPass:   viper.GetString("DB_PASS"),
		DBName:   viper.GetString("DB_NAME"),
		DBHost:   viper.GetString("DB_HOST"),
		DBPort:   viper.GetInt("DB_PORT"),
		APPPort:  viper.GetInt("APP_PORT"),
	}
}

var config = load()

func Cfg() *Config {
	return config
}
