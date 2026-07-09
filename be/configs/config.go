package configs

import "github.com/spf13/viper"

type Config struct {
	App AppConfig `json:"app"`
}

type AppConfig struct {
	Env  string `json:"APP_ENV"`
	Port int    `json:"APP_PORT"`
}

func NewConfig() *Config {
	return &Config{
		App: AppConfig{
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetInt("APP_PORT"),
		},
	}
}
