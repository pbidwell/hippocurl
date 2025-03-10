package utils

// ContextKey is a custom type to avoid key collisions in context
type ContextKey string

const (
	LoggerKey         ContextKey = "logger"
	LogFilePath       ContextKey = "logfilepath"
	ConfigKey         ContextKey = "config"
	ConfigFilePathKey ContextKey = "configfilepath"
)
