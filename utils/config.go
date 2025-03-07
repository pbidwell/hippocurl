package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	HC_CONFIG_DIR_NAME = ".hcconfig"
	HC_CONFIG_FILE     = "hc_config.yml"
)

// Config represents the structure of the YAML configuration file
type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name         string        `yaml:"name"`
	Environments []Environment `yaml:"environments"`
	Routes       []Route       `yaml:"routes"`
}

type Environment struct {
	Name    string `yaml:"name"`
	BaseURL string `yaml:"base_url"`
	Auth    Auth   `yaml:"auth"`
}

type Auth struct {
	Type     string `yaml:"type"` // e.g., "basic", "bearer", "none"
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token,omitempty"`
}

type Route struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Method      string `yaml:"method"`
	Description string `yaml:"description"`
}

// LoadConfig reads the YAML configuration file and stores it in the context
func LoadConfig(ctx context.Context) context.Context {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get user home directory: %v", err)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, HC_CONFIG_DIR_NAME)
	configFilePath := filepath.Join(configDir, HC_CONFIG_FILE)

	// Ensure the config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("failed to create config directory: %v", err)
		os.Exit(1)
	}

	// Check if the config file exists, if not, create an empty config file
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		defaultConfig := &Config{Services: []Service{}}
		data, _ := yaml.Marshal(defaultConfig)
		if err := ioutil.WriteFile(configFilePath, data, 0644); err != nil {
			fmt.Printf("failed to create default config file: %v", err)
			os.Exit(1)
		}
	}

	// Read the config file
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("failed to read config file: %v", err)
		os.Exit(1)
	}

	// Unmarshal the YAML data into the Config struct
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("failed to parse config file: %v", err)
		os.Exit(1)
	}

	// Store config in context
	ctx = context.WithValue(ctx, ConfigKey, &config)
	ctx = context.WithValue(ctx, ConfigFilePathKey, configFilePath)
	return ctx
}

func (c Config) GetServiceNames() []string {
	var names []string
	for _, service := range c.Services {
		names = append(names, service.Name)
	}
	return names
}

func (c Config) GetServiceByName(name string) *Service {
	for _, service := range c.Services {
		if service.Name == name {
			return &service
		}
	}
	return nil
}

func (s Service) GetRouteNames() []string {
	var names []string
	for _, route := range s.Routes {
		names = append(names, route.Name)
	}
	return names
}

func (s Service) GetRouteByName(name string) *Route {
	for _, route := range s.Routes {
		if route.Name == name {
			return &route
		}
	}
	return nil
}

func (s Service) GetEnvironmentNames() []string {
	var names []string
	for _, env := range s.Environments {
		names = append(names, env.Name)
	}
	return names
}

func (s Service) GetEnvironmentByName(name string) *Environment {
	for _, env := range s.Environments {
		if env.Name == name {
			return &env
		}
	}
	return nil
}
