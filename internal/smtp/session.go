package smtp

import (
	"context"
	"io"
	"log/slog"

	"github.com/alex/smtp-gotify/internal/mail"
	"github.com/emersion/go-smtp"
)

type Session struct {
	logger    *slog.Logger
	parser    *mail.Parser
	forwarder Forwarder
	from      string
	to        []string
}

func NewSession(logger *slog.Logger, parser *mail.Parser, forwarder Forwarder) *Session {
	return &Session{
		logger:    logger,
		parser:    parser,
		forwarder: forwarder,
	}
}

func (s *Session) AuthPlain(username, password string) error {
	// Accept any authentication
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	s.logger.Debug("mail from", "from", from)
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = append(s.to, to)
	s.logger.Debug("rcpt to", "to", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	msg, err := s.parser.Parse(r)
	if err != nil {
		s.logger.Error("failed to parse email", "error", err)
		return err
	}

	// Use envelope addresses if headers are empty
	if msg.From == "" {
		msg.From = s.from
	}
	if len(msg.To) == 0 {
		msg.To = s.to
	}

	s.logger.Info("received email",
		"from", msg.From,
		"to", msg.To,
		"subject", msg.Subject,
		"attachments", len(msg.Attachments),
	)

	ctx := context.Background()
	if err := s.forwarder.Forward(ctx, msg); err != nil {
		s.logger.Error("failed to forward to gotify", "error", err)
		return err
	}

	s.logger.Info("forwarded to gotify", "subject", msg.Subject)
	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.to = nil
}

func (s *Session) Logout() error {
	return nil
}
