package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	hcConfigDirName     = ".hc"
	hcLogFileName       = "hc.log"
	hcAPIConfigFileName = "api_config.yml"
)

func Load() *App {
	configDir := createHCDirectory()

	logger, logFilePath := buildLogger(configDir)
	return &App{
		Logger:      logger,
		LogFilePath: logFilePath,
		APIConfig:   loadAPIConfig(configDir),
	}
}

func createHCDirectory() string {
	// Expand home directory if necessary
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to detect home directory: %v\n", err)
	}

	configDir := filepath.Join(homeDir, hcConfigDirName)

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v\n", err)
	}

	return configDir
}

func buildLogger(configDir string) (*log.Logger, string) {
	logFilePath := filepath.Join(configDir, hcLogFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v\n", err)
	}
	return log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile), logFilePath
}

func loadAPIConfig(configDir string) *APIConfig {
	path := filepath.Join(configDir, hcAPIConfigFileName)

	v := viper.New()
	v.SetConfigFile(path)   // path to config file (e.g., "./config.yml")
	v.SetConfigType("yaml") // optional if file extension is clear

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("read config: %v", err)
	}

	var cfg APIConfig
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("unmarshal config: %v", err)
	}

	return &cfg
}
