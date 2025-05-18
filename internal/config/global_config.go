/*
Copyright © 2025 Pablo Bidwell <bidwell.pablo@gmail.com>
*/
package config

import "log"

type App struct {
	GlobalConfig *GlobalConfig
	APIConfig    *APIConfig
	Logger       *log.Logger
	LogFilePath  string
}
type GlobalConfig struct {
	// global hc file configuration
	FilePath string
}
