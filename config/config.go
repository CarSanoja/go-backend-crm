package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string
	CSVFile string
}

var config Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	config = Config{
		Port:    viper.GetString("port"),
		CSVFile: viper.GetString("csv_file"),
	}
}

func GetConfig() Config {
	return config
}
