package config

import "os"

type Config struct {
	QWeatherAPIKey  string
	QWeatherBaseURL string
}

func NewConfig() *Config {
	return &Config{
		QWeatherAPIKey:  getEnv("QWEATHER_API_KEY", ""),
		QWeatherBaseURL: getEnv("QWEATHER_BASE_URL", ""), // https://xxx
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
