package smtp

import (
	"context"
	"log/slog"

	"github.com/alex/smtp-gotify/internal/mail"
	"github.com/emersion/go-smtp"
)

type Forwarder interface {
	Forward(ctx context.Context, msg *mail.Message) error
}

type Backend struct {
	logger    *slog.Logger
	parser    *mail.Parser
	forwarder Forwarder
}

func NewBackend(logger *slog.Logger, parser *mail.Parser, forwarder Forwarder) *Backend {
	return &Backend{
		logger:    logger,
		parser:    parser,
		forwarder: forwarder,
	}
}

func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	if c != nil {
		b.logger.Debug("new session", "remote", c.Conn().RemoteAddr())
	}
	return NewSession(b.logger, b.parser, b.forwarder), nil
}
