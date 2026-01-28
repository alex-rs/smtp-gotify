package smtp

import (
	"log/slog"
	"time"

	"github.com/alex/smtp-gotify/internal/config"
	"github.com/alex/smtp-gotify/internal/mail"
	"github.com/emersion/go-smtp"
)

type Server struct {
	server *smtp.Server
	logger *slog.Logger
}

func NewServer(cfg config.SMTPConfig, logger *slog.Logger, parser *mail.Parser, forwarder Forwarder) *Server {
	backend := NewBackend(logger, parser, forwarder)

	s := smtp.NewServer(backend)
	s.Addr = cfg.Listen
	s.Domain = cfg.Domain
	s.MaxMessageBytes = int64(cfg.MaxSize)
	s.ReadTimeout = 60 * time.Second
	s.WriteTimeout = 60 * time.Second
	s.AllowInsecureAuth = true

	return &Server{
		server: s,
		logger: logger,
	}
}

func (s *Server) ListenAndServe() error {
	s.logger.Info("starting SMTP server", "addr", s.server.Addr, "domain", s.server.Domain)
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	s.logger.Info("stopping SMTP server")
	return s.server.Close()
}
