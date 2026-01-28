package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set required env vars
	os.Setenv("GOTIFY_URL", "http://localhost:8080")
	os.Setenv("GOTIFY_TOKEN", "test-token")
	defer func() {
		os.Unsetenv("GOTIFY_URL")
		os.Unsetenv("GOTIFY_TOKEN")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Gotify.URL != "http://localhost:8080" {
		t.Errorf("expected URL http://localhost:8080, got %s", cfg.Gotify.URL)
	}

	if len(cfg.Gotify.Tokens) != 1 || cfg.Gotify.Tokens[0] != "test-token" {
		t.Errorf("expected token test-token, got %v", cfg.Gotify.Tokens)
	}

	// Check defaults
	if cfg.Gotify.Priority != 5 {
		t.Errorf("expected priority 5, got %d", cfg.Gotify.Priority)
	}

	if cfg.SMTP.Listen != ":2525" {
		t.Errorf("expected SMTP listen :2525, got %s", cfg.SMTP.Listen)
	}

	if cfg.Health.Enabled != true {
		t.Errorf("expected health enabled true, got %v", cfg.Health.Enabled)
	}
}

func TestLoadMissingRequired(t *testing.T) {
	os.Unsetenv("GOTIFY_URL")
	os.Unsetenv("GOTIFY_TOKEN")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing required vars")
	}
}

func TestLoadMultipleTokens(t *testing.T) {
	os.Setenv("GOTIFY_URL", "http://localhost:8080")
	os.Setenv("GOTIFY_TOKEN", "token1, token2, token3")
	defer func() {
		os.Unsetenv("GOTIFY_URL")
		os.Unsetenv("GOTIFY_TOKEN")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Gotify.Tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(cfg.Gotify.Tokens))
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid",
			cfg: Config{
				Gotify: GotifyConfig{URL: "http://example.com", Tokens: []string{"tok"}, Priority: 5},
				SMTP:   SMTPConfig{MaxSize: 1000},
				Log:    LogConfig{Level: "info", Format: "json"},
			},
			wantErr: false,
		},
		{
			name: "invalid priority",
			cfg: Config{
				Gotify: GotifyConfig{URL: "http://example.com", Tokens: []string{"tok"}, Priority: 15},
				SMTP:   SMTPConfig{MaxSize: 1000},
				Log:    LogConfig{Level: "info", Format: "json"},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			cfg: Config{
				Gotify: GotifyConfig{URL: "http://example.com", Tokens: []string{"tok"}, Priority: 5},
				SMTP:   SMTPConfig{MaxSize: 1000},
				Log:    LogConfig{Level: "invalid", Format: "json"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
