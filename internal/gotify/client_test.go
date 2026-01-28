package gotify

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/smtp-gotify/internal/mail"
	"github.com/alex/smtp-gotify/internal/template"
)

func TestClient_Forward(t *testing.T) {
	var received Message
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Query().Get("token") != "test-token" {
			t.Errorf("expected token test-token, got %s", r.URL.Query().Get("token"))
		}

		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	renderer, _ := template.NewRenderer("{{.Subject}}", "{{.Body}}")
	client := NewClient(Config{
		URL:      server.URL,
		Tokens:   []string{"test-token"},
		Priority: 5,
		Renderer: renderer,
		Logger:   slog.Default(),
	})

	msg := &mail.Message{
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err := client.Forward(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Title != "Test Subject" {
		t.Errorf("expected title 'Test Subject', got %s", received.Title)
	}

	if received.Priority != 5 {
		t.Errorf("expected priority 5, got %d", received.Priority)
	}
}

func TestClient_ForwardWithMarkdown(t *testing.T) {
	var received Message
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	renderer, _ := template.NewRenderer("{{.Subject}}", "{{.Body}}")
	client := NewClient(Config{
		URL:      server.URL,
		Tokens:   []string{"test-token"},
		Priority: 5,
		Markdown: true,
		Renderer: renderer,
		Logger:   slog.Default(),
	})

	msg := &mail.Message{Subject: "Test", Body: "Body"}
	err := client.Forward(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Extras == nil {
		t.Fatal("expected extras for markdown")
	}

	display, ok := received.Extras["client::display"].(map[string]interface{})
	if !ok {
		t.Fatal("expected client::display in extras")
	}

	if display["contentType"] != "text/markdown" {
		t.Errorf("expected contentType text/markdown, got %v", display["contentType"])
	}
}

func TestClient_ForwardMultipleTokens(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	renderer, _ := template.NewRenderer("{{.Subject}}", "{{.Body}}")
	client := NewClient(Config{
		URL:      server.URL,
		Tokens:   []string{"token1", "token2", "token3"},
		Renderer: renderer,
		Logger:   slog.Default(),
	})

	msg := &mail.Message{Subject: "Test", Body: "Body"}
	err := client.Forward(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestClient_ForwardError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	renderer, _ := template.NewRenderer("{{.Subject}}", "{{.Body}}")
	client := NewClient(Config{
		URL:      server.URL,
		Tokens:   []string{"test-token"},
		Renderer: renderer,
		Logger:   slog.Default(),
	})

	msg := &mail.Message{Subject: "Test", Body: "Body"}
	err := client.Forward(context.Background(), msg)
	if err == nil {
		t.Error("expected error for server error response")
	}
}
