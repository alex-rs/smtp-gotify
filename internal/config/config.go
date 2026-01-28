package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Gotify GotifyConfig
	SMTP   SMTPConfig
	Health HealthConfig
	Log    LogConfig
}

type GotifyConfig struct {
	URL             string
	Tokens          []string
	Priority        int
	Markdown        bool
	TitleTemplate   string
	MessageTemplate string
}

type SMTPConfig struct {
	Listen  string
	Domain  string
	MaxSize int
}

type HealthConfig struct {
	Enabled bool
	Listen  string
}

type LogConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	cfg := &Config{
		Gotify: GotifyConfig{
			URL:             getEnv("GOTIFY_URL", ""),
			Tokens:          parseTokens(getEnv("GOTIFY_TOKEN", "")),
			Priority:        getEnvInt("GOTIFY_PRIORITY", 5),
			Markdown:        getEnvBool("GOTIFY_MARKDOWN", false),
			TitleTemplate:   getEnv("GOTIFY_TITLE_TEMPLATE", "{{.Subject}}"),
			MessageTemplate: getEnv("GOTIFY_MESSAGE_TEMPLATE", "From: {{.From}}\nTo: {{.To}}\n---\n{{.Body}}"),
		},
		SMTP: SMTPConfig{
			Listen:  getEnv("SMTP_LISTEN", ":2525"),
			Domain:  getEnv("SMTP_DOMAIN", "localhost"),
			MaxSize: getEnvInt("SMTP_MAX_SIZE", 10485760),
		},
		Health: HealthConfig{
			Enabled: getEnvBool("HEALTH_ENABLED", true),
			Listen:  getEnv("HEALTH_LISTEN", ":8080"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var errs []error

	if c.Gotify.URL == "" {
		errs = append(errs, errors.New("GOTIFY_URL is required"))
	}

	if len(c.Gotify.Tokens) == 0 {
		errs = append(errs, errors.New("GOTIFY_TOKEN is required"))
	}

	if c.Gotify.Priority < 0 || c.Gotify.Priority > 10 {
		errs = append(errs, fmt.Errorf("GOTIFY_PRIORITY must be between 0 and 10, got %d", c.Gotify.Priority))
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Log.Level] {
		errs = append(errs, fmt.Errorf("LOG_LEVEL must be one of debug/info/warn/error, got %s", c.Log.Level))
	}

	validFormats := map[string]bool{"json": true, "text": true}
	if !validFormats[c.Log.Format] {
		errs = append(errs, fmt.Errorf("LOG_FORMAT must be one of json/text, got %s", c.Log.Format))
	}

	if c.SMTP.MaxSize <= 0 {
		errs = append(errs, fmt.Errorf("SMTP_MAX_SIZE must be positive, got %d", c.SMTP.MaxSize))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

func parseTokens(s string) []string {
	if s == "" {
		return nil
	}
	var tokens []string
	for _, t := range strings.Split(s, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	return tokens
}
