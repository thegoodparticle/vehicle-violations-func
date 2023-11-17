package config

import (
	"os"
	"strconv"
)

func GetEnv(env, defaultValue string) string {
	environment := os.Getenv(env)
	if environment == "" {
		return defaultValue
	}

	return environment
}

type Config struct {
	Port    int
	Timeout int
}

func GetConfig() Config {
	return Config{
		Port:    parseEnvToInt("PORT", "80"),
		Timeout: parseEnvToInt("TIMEOUT", "30"),
	}
}

func parseEnvToInt(envName, defaultValue string) int {
	num, err := strconv.Atoi(GetEnv(envName, defaultValue))

	if err != nil {
		return 0
	}

	return num
}
