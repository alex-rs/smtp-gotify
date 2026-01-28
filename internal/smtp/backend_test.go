package smtp

import (
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/alex/smtp-gotify/internal/mail"
)

type mockForwarder struct {
	messages []*mail.Message
	err      error
}

func (m *mockForwarder) Forward(ctx context.Context, msg *mail.Message) error {
	m.messages = append(m.messages, msg)
	return m.err
}

func TestSession_Data(t *testing.T) {
	forwarder := &mockForwarder{}
	parser := mail.NewParser()
	session := NewSession(slog.Default(), parser, forwarder)

	session.Mail("sender@example.com", nil)
	session.Rcpt("recipient@example.com", nil)

	email := `From: sender@example.com
To: recipient@example.com
Subject: Test
Content-Type: text/plain

Hello world`

	err := session.Data(strings.NewReader(email))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(forwarder.messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(forwarder.messages))
	}

	msg := forwarder.messages[0]
	if msg.Subject != "Test" {
		t.Errorf("expected subject Test, got %s", msg.Subject)
	}
}

func TestSession_Reset(t *testing.T) {
	forwarder := &mockForwarder{}
	parser := mail.NewParser()
	session := NewSession(slog.Default(), parser, forwarder)

	session.Mail("sender@example.com", nil)
	session.Rcpt("recipient@example.com", nil)
	session.Reset()

	if session.from != "" {
		t.Errorf("expected from to be empty after reset")
	}

	if len(session.to) != 0 {
		t.Errorf("expected to to be empty after reset")
	}
}

func TestSession_EnvelopeAddresses(t *testing.T) {
	forwarder := &mockForwarder{}
	parser := mail.NewParser()
	session := NewSession(slog.Default(), parser, forwarder)

	session.Mail("envelope-sender@example.com", nil)
	session.Rcpt("envelope-recipient@example.com", nil)

	// Email without From/To headers
	email := `Subject: No Headers
Content-Type: text/plain

Body only`

	err := session.Data(strings.NewReader(email))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := forwarder.messages[0]
	if msg.From != "envelope-sender@example.com" {
		t.Errorf("expected envelope from, got %s", msg.From)
	}

	if len(msg.To) != 1 || msg.To[0] != "envelope-recipient@example.com" {
		t.Errorf("expected envelope to, got %v", msg.To)
	}
}

func TestBackend_NewSession(t *testing.T) {
	forwarder := &mockForwarder{}
	parser := mail.NewParser()
	backend := NewBackend(slog.Default(), parser, forwarder)

	session, err := backend.NewSession(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session == nil {
		t.Error("expected session to be created")
	}
}
