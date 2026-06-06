package config_test

import (
	"testing"

	"github.com/adamkirk/panoptes/api/internal/config"
)

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}
}

func TestGetServerPort(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{Port: 9000}}
	if cfg.GetServerPort() != 9000 {
		t.Errorf("expected 9000, got %d", cfg.GetServerPort())
	}
}
