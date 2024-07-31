package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	UploadDir string
	JWTSecret string
}

var config Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer el archivo de configuración: %v", err)
	}

	config = Config{
		Port:      viper.GetString("port"),
		UploadDir: viper.GetString("upload_dir"),
		JWTSecret: viper.GetString("jwt_secret"),
	}
}

func GetConfig() Config {
	return config
}
