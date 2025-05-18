/*
Copyright © 2025 Pablo Bidwell <bidwell.pablo@gmail.com>
*/
package config

type APIConfig struct {
	Services []Service `mapstructure:"services"`
}

type Service struct {
	Name         string        `mapstructure:"name"`
	Environments []Environment `mapstructure:"environments"`
	Routes       []Route       `mapstructure:"routes"`
}

type Environment struct {
	Name    string            `mapstructure:"name"`
	BaseURL string            `mapstructure:"base_url"`
	Auth    Auth              `mapstructure:"auth"`
	Headers map[string]string `mapstructure:"headers,omitempty"` // Custom headers
}

type Auth struct {
	Type     string `mapstructure:"type"` // e.g., "basic", "bearer", "none"
	Username string `mapstructure:"username,omitempty"`
	Password string `mapstructure:"password,omitempty"`
	Token    string `mapstructure:"token,omitempty"`
}

type Route struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Method      string `mapstructure:"method"`
	Path        string `mapstructure:"path"`
	Body        string `mapstructure:"body"`
}

// Allows adherence to the Named interface
func (s Service) GetName() string     { return s.Name }
func (r Route) GetName() string       { return r.Name }
func (e Environment) GetName() string { return e.Name }

// Represents config elements with Name fields
type Named interface {
	GetName() string
}

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

func (c *APIConfig) GetServiceNames() []string {
	return getNames(c.Services)
}

func (c *APIConfig) GetServiceByName(name string) *Service {
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
