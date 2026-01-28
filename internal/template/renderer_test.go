package template

import (
	"testing"

	"github.com/alex/smtp-gotify/internal/mail"
)

func TestRenderer_Render(t *testing.T) {
	r, err := NewRenderer("{{.Subject}}", "From: {{.From}}\nTo: {{.To}}\n---\n{{.Body}}")
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	msg := &mail.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test body content",
	}

	title, body, err := r.Render(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if title != "Test Subject" {
		t.Errorf("expected title 'Test Subject', got %s", title)
	}

	expectedBody := "From: sender@example.com\nTo: recipient@example.com\n---\nTest body content"
	if body != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}
}

func TestRenderer_RenderMultipleRecipients(t *testing.T) {
	r, err := NewRenderer("{{.Subject}}", "To: {{.To}}")
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	msg := &mail.Message{
		To: []string{"a@example.com", "b@example.com", "c@example.com"},
	}

	_, body, err := r.Render(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "To: a@example.com, b@example.com, c@example.com"
	if body != expected {
		t.Errorf("expected %q, got %q", expected, body)
	}
}

func TestNewRenderer_InvalidTemplate(t *testing.T) {
	_, err := NewRenderer("{{.Invalid", "valid")
	if err == nil {
		t.Error("expected error for invalid template")
	}
}
