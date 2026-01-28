package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alex/smtp-gotify/internal/config"
	"github.com/alex/smtp-gotify/internal/gotify"
	"github.com/alex/smtp-gotify/internal/health"
	"github.com/alex/smtp-gotify/internal/mail"
	"github.com/alex/smtp-gotify/internal/smtp"
	"github.com/alex/smtp-gotify/internal/template"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := setupLogger(cfg.Log)

	renderer, err := template.NewRenderer(cfg.Gotify.TitleTemplate, cfg.Gotify.MessageTemplate)
	if err != nil {
		logger.Error("failed to create template renderer", "error", err)
		os.Exit(1)
	}

	parser := mail.NewParser()

	gotifyClient := gotify.NewClient(gotify.Config{
		URL:      cfg.Gotify.URL,
		Tokens:   cfg.Gotify.Tokens,
		Priority: cfg.Gotify.Priority,
		Markdown: cfg.Gotify.Markdown,
		Renderer: renderer,
		Logger:   logger,
	})

	smtpServer := smtp.NewServer(cfg.SMTP, logger, parser, gotifyClient)

	var healthServer *health.Server
	if cfg.Health.Enabled {
		healthServer = health.NewServer(cfg.Health.Listen, logger)
		go func() {
			if err := healthServer.ListenAndServe(); err != nil {
				logger.Error("health server error", "error", err)
			}
		}()
	}

	go func() {
		if err := smtpServer.ListenAndServe(); err != nil {
			logger.Error("SMTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutting down")

	if healthServer != nil {
		if err := healthServer.Shutdown(context.Background()); err != nil {
			logger.Error("health server shutdown error", "error", err)
		}
	}
	if err := smtpServer.Close(); err != nil {
		logger.Error("SMTP server close error", "error", err)
	}
}

func setupLogger(cfg config.LogConfig) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if cfg.Format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
