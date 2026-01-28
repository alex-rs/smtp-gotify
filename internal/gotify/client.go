package gotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/alex/smtp-gotify/internal/mail"
	"github.com/alex/smtp-gotify/internal/template"
)

type Message struct {
	Title    string                 `json:"title"`
	Message  string                 `json:"message"`
	Priority int                    `json:"priority"`
	Extras   map[string]interface{} `json:"extras,omitempty"`
}

type Client struct {
	baseURL  string
	tokens   []string
	priority int
	markdown bool
	renderer *template.Renderer
	http     *http.Client
	logger   *slog.Logger
}

type Config struct {
	URL      string
	Tokens   []string
	Priority int
	Markdown bool
	Renderer *template.Renderer
	Logger   *slog.Logger
}

func NewClient(cfg Config) *Client {
	return &Client{
		baseURL:  strings.TrimSuffix(cfg.URL, "/"),
		tokens:   cfg.Tokens,
		priority: cfg.Priority,
		markdown: cfg.Markdown,
		renderer: cfg.Renderer,
		logger:   cfg.Logger,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Forward(ctx context.Context, msg *mail.Message) error {
	title, body, err := c.renderer.Render(msg)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	gotifyMsg := Message{
		Title:    title,
		Message:  body,
		Priority: c.priority,
	}

	if c.markdown {
		gotifyMsg.Extras = map[string]interface{}{
			"client::display": map[string]string{
				"contentType": "text/markdown",
			},
		}
	}

	payload, err := json.Marshal(gotifyMsg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	// Send to all configured tokens
	var errs []error
	for _, token := range c.tokens {
		if err := c.send(ctx, token, payload); err != nil {
			errs = append(errs, fmt.Errorf("token %s...: %w", token[:min(8, len(token))], err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send to %d/%d tokens: %v", len(errs), len(c.tokens), errs)
	}

	return nil
}

func (c *Client) send(ctx context.Context, token string, payload []byte) error {
	url := fmt.Sprintf("%s/message?token=%s", c.baseURL, token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	c.logger.Debug("message sent to gotify", "status", resp.Status)
	return nil
}
