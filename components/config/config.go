package config

import (
	"os"
)

type AppConfig struct {
	Env		  string
	Port      string
	TableName string
	Region    string
	BaseURL   string
}

func LoadConfig() *AppConfig {
	cfg := &AppConfig{
		Env:	   getEnv("ENV", "dev"),
		Port:      getEnv("PORT", "80"),
		TableName: getEnv("DDB_TABLE", "url"),
		Region:    getEnv("AWS_REGION", "us-east-1"),
		BaseURL:   getEnv("BASE_URL", "http://localhost:8080/"),
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
