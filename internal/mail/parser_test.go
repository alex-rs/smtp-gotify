package mail

import (
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	p := NewParser()

	email := `From: sender@example.com
To: recipient@example.com
Subject: Test Subject
Content-Type: text/plain

This is the body.`

	msg, err := p.Parse(strings.NewReader(email))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.From != "sender@example.com" {
		t.Errorf("expected from sender@example.com, got %s", msg.From)
	}

	if msg.Subject != "Test Subject" {
		t.Errorf("expected subject Test Subject, got %s", msg.Subject)
	}

	if !strings.Contains(msg.Body, "This is the body") {
		t.Errorf("expected body to contain 'This is the body', got %s", msg.Body)
	}
}

func TestParser_ParseMultipart(t *testing.T) {
	p := NewParser()

	email := `From: sender@example.com
To: recipient@example.com
Subject: Multipart Test
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary"

--boundary
Content-Type: text/plain

Plain text body.
--boundary
Content-Type: text/html

<p>HTML body</p>
--boundary--`

	msg, err := p.Parse(strings.NewReader(email))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(msg.Body, "Plain text body") {
		t.Errorf("expected plain text body, got %s", msg.Body)
	}
}

func TestParser_ParseInvalid(t *testing.T) {
	p := NewParser()

	// Empty input should return an empty message (enmime is lenient)
	msg, err := p.Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.Subject != "" {
		t.Errorf("expected empty subject, got %s", msg.Subject)
	}
}
