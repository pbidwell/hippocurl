package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ContextKey is a custom type to avoid key collisions in context
type ContextKey string

const (
	LoggerKey         ContextKey = "logger"
	LogFilePath       ContextKey = "logfilepath"
	ConfigKey         ContextKey = "config"
	ConfigFilePathKey ContextKey = "configfilepath"

	HC_CONFIG_DIR = ".hcconfig"
	HC_LOG_FILE   = "hc.log"
)

func LoadLoggerIntoContext(ctx context.Context) context.Context {
	// Expand home directory if necessary
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(homeDir, HC_CONFIG_DIR)
		logFilePath := filepath.Join(configDir, HC_LOG_FILE)

		// Ensure config directory exists
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("Failed to create config directory: %v\n", err)
		}

		// Setup logging
		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open log file: %v\n", err)
			os.Exit(1)
		}
		logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

		// Store logger and config in context
		ctx = context.WithValue(ctx, LoggerKey, logger)
		ctx = context.WithValue(ctx, LogFilePath, logFilePath)
	}
	return ctx
}
