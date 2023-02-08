package app

import (
	"bufio"
	"os"
	"strings"
)

type OAuthConfig struct {
	ClientId     string
	ClientSecret string
}

type Config struct {
	Google OAuthConfig
	GitHub OAuthConfig
}

func NewConfig() *Config {
	return &Config{
		Google: OAuthConfig{
			ClientId:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		},
		GitHub: OAuthConfig{
			ClientId:     getEnv("GITHUB_CLIENT_ID", ""),
			ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func SetEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) == 2 {
			key, value := strings.TrimSpace(line[0]), strings.Trim(strings.TrimSpace(line[1]), "\"")
			if key != "" && value != "" {
				os.Setenv(key, value)
			}
		}
	}
	return scanner.Err()
}
