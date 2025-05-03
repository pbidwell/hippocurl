package utils

import (
	"context"
	"fmt"
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
	Name    string            `yaml:"name"`
	BaseURL string            `yaml:"base_url"`
	Auth    Auth              `yaml:"auth"`
	Headers map[string]string `yaml:"headers,omitempty"` // Custom headers
}

type Auth struct {
	Type     string `yaml:"type"` // e.g., "basic", "bearer", "none"
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token,omitempty"`
}

type Route struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Method      string `yaml:"method"`
	Path        string `yaml:"path"`
	Body        string `yaml:"body"`
}

// Represents config elements with Name fields
type Named interface {
	GetName() string
}

// Allows adherence to the Named interface
func (s Service) GetName() string     { return s.Name }
func (r Route) GetName() string       { return r.Name }
func (e Environment) GetName() string { return e.Name }

// getNames acts as a generic helper function to avoid
// duplicate code on "get<ENTITY>Names()" calls
func getNames[T Named](items []T) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.GetName()
	}
	return names
}

// getNames acts as a generic helper function to avoid
// duplicate code on "get<ENTITY>ByName()" calls
func getByName[T Named](items []T, name string) *T {
	for i := range items {
		if items[i].GetName() == name {
			return &items[i]
		}
	}
	return nil
}

func (c *Config) GetServiceNames() []string {
	return getNames(c.Services)
}

func (c *Config) GetServiceByName(name string) *Service {
	return getByName(c.Services, name)
}

func (s *Service) GetRouteNames() []string {
	return getNames(s.Routes)
}

func (s *Service) GetRouteByName(name string) *Route {
	return getByName(s.Routes, name)
}

func (s *Service) GetEnvironmentNames() []string {
	return getNames(s.Environments)
}

func (s *Service) GetEnvironmentByName(name string) *Environment {
	return getByName(s.Environments, name)
}

func loadConfigFromFilepath(path string) (*Config, error) {
	// Read the config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("config file at %s is empty", path)
	}

	// Unmarshal the YAML data into the Config struct
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// LoadConfigIntoContext reads the YAML configuration file and stores it in the context
func LoadConfigIntoContext(ctx context.Context) (context.Context, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ctx, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, HC_CONFIG_DIR_NAME)
	configFilePath := filepath.Join(configDir, HC_CONFIG_FILE)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return ctx, fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	// Create default config if missing
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		defaultConfig := &Config{Services: []Service{}}
		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return ctx, fmt.Errorf("failed to marshal default config: %w", err)
		}
		if err := os.WriteFile(configFilePath, data, 0644); err != nil {
			return ctx, fmt.Errorf("failed to write default config to %s: %w", configFilePath, err)
		}
	}

	config, err := loadConfigFromFilepath(configFilePath)
	if err != nil {
		return ctx, fmt.Errorf("failed to load config from %s: %w", configFilePath, err)
	}

	ctx = context.WithValue(ctx, ConfigKey, config)
	ctx = context.WithValue(ctx, ConfigFilePathKey, configFilePath)
	return ctx, nil
}
