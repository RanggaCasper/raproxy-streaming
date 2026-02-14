package config

import (
	"time"
)

// Config holds application configuration
type Config struct {
	HTTP HTTPConfig
}

// HTTPConfig holds HTTP client configuration
type HTTPConfig struct {
	Timeout        time.Duration
	ConnectTimeout time.Duration
	MaxRedirects   int
}

// New creates a new configuration with default values
func New() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Timeout:        60 * time.Second,
			ConnectTimeout: 10 * time.Second,
			MaxRedirects:   10,
		},
	}
}
