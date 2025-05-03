package utils

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	minimalConfig *Config
	normalConfig  *Config
	emptyConfig   *Config
)

func TestMain(m *testing.M) {
	var err error
	base := "testdata"

	// Load configs once
	minimalConfig, err = loadConfigFromFilepath(filepath.Join(base, "minimal_config.yml"))
	if err != nil {
		panic(err)
	}

	normalConfig, err = loadConfigFromFilepath(filepath.Join(base, "normal_config.yml"))
	if err != nil {
		panic(err)
	}

	emptyConfig, err = loadConfigFromFilepath(filepath.Join(base, "empty_config.yml"))
	if err != nil {
		panic(err)
	}

	// Run tests
	os.Exit(m.Run())
}

func TestGetServiceNames(t *testing.T) {
	// empty config
	serviceNames := emptyConfig.GetServiceNames()
	if len(serviceNames) != 0 {
		t.Errorf("unexpected service names from empty config: %v", serviceNames)
	}

	// minimal config
	serviceNames = minimalConfig.GetServiceNames()
	if len(serviceNames) != 1 || serviceNames[0] != "ServiceOne" {
		t.Errorf("unexpected service names from minimal config: %v", serviceNames)
	}
}

func TestGetServiceByName(t *testing.T) {
	// empty config
	service := emptyConfig.GetServiceByName("nonexistent")
	if service != nil {
		t.Errorf("unexpected service entity extraction from test config")
	}

	// minimal config
	service = minimalConfig.GetServiceByName("ServiceOne")
	if service == nil {
		t.Errorf("failed to extract service from config by name")
	}
}

func TestGetRouteNames(t *testing.T) {
	service := minimalConfig.GetServiceByName("ServiceOne")
	if service == nil {
		t.Fatal("ServiceOne not found in minimalConfig")
	}

	routes := service.GetRouteNames()
	if len(routes) != 1 || routes[0] != "RouteNameOne" {
		t.Errorf("expected 1 route named RouteNameOne, got: %v", routes)
	}
}

func TestGetRouteByName(t *testing.T) {
	service := minimalConfig.GetServiceByName("ServiceOne")
	if service == nil {
		t.Fatal("ServiceOne not found in minimalConfig")
	}

	// Nonexistent route
	if route := service.GetRouteByName("nonexistent"); route != nil {
		t.Errorf("expected nil for nonexistent route")
	}

	// Existing route
	route := service.GetRouteByName("RouteNameOne")
	if route == nil {
		t.Errorf("expected to find RouteNameOne")
	}
}

func TestGetEnvironmentNames(t *testing.T) {
	service := minimalConfig.GetServiceByName("ServiceOne")
	if service == nil {
		t.Fatal("ServiceOne not found in minimalConfig")
	}

	envs := service.GetEnvironmentNames()
	if len(envs) != 1 || envs[0] != "NameOne" {
		t.Errorf("expected 1 environment named NameOne, got: %v", envs)
	}
}

func TestGetEnvironmentByName(t *testing.T) {
	service := minimalConfig.GetServiceByName("ServiceOne")
	if service == nil {
		t.Fatal("ServiceOne not found in minimalConfig")
	}

	// Nonexistent environment
	if env := service.GetEnvironmentByName("nonexistent"); env != nil {
		t.Errorf("expected nil for nonexistent environment")
	}

	// Existing environment
	env := service.GetEnvironmentByName("NameOne")
	if env == nil {
		t.Errorf("expected to find environment NameOne")
	}
}
