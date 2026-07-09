package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig `json:"app"`
	Database Database  `json:"database"`
}

type AppConfig struct {
	Env  string `json:"APP_ENV"`
	Port int    `json:"APP_PORT"`
}

type Database struct {
	Host           string   `json:"host"`
	Port           int      `json:"port"`
	User           string   `json:"user"`
	Pass           string   `json:"pass"`
	Name           string   `json:"name"`
	Schema         string   `json:"schema"`
	SSLMode        string   `json:"ssl_mode"`
	SlaveHosts     []string `json:"slave_hosts"`
	DBMaxOpenConns int      `json:"db_max_open_conns"`
	DBMaxIdleConns int      `json:"db_max_idle_conns"`
}

func NewConfig() *Config {
	return &Config{
		App: AppConfig{
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetInt("APP_PORT"),
		},
		Database: Database{
			Host:           viper.GetString("DB_HOST"),
			Port:           viper.GetInt("DB_PORT"),
			User:           viper.GetString("DB_USER"),
			Pass:           viper.GetString("DB_PASS"),
			Name:           viper.GetString("DB_NAME"),
			Schema:         viper.GetString("DB_SCHEMA"),
			SSLMode:        getStringOrDefault(viper.GetString("DB_SSL_MODE"), "disable"),
			SlaveHosts:     splitCSV(viper.GetString("DB_SLAVE_HOSTS")),
			DBMaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONN"),
			DBMaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONN"),
		},
	}
}

func splitCSV(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	rawParts := strings.Split(input, ",")
	parts := make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}

	return parts
}

func getStringOrDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}
