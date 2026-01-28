package mail

import (
	"io"

	"github.com/jhillyerd/enmime"
)

type Message struct {
	From        string
	To          []string
	Subject     string
	Body        string
	Attachments []Attachment
}

type Attachment struct {
	Filename    string
	ContentType string
	Size        int
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader) (*Message, error) {
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		From:    env.GetHeader("From"),
		Subject: env.GetHeader("Subject"),
		To:      parseAddressList(env),
	}

	// Prefer plain text, fall back to HTML
	if env.Text != "" {
		msg.Body = env.Text
	} else if env.HTML != "" {
		msg.Body = env.HTML
	}

	for _, att := range env.Attachments {
		msg.Attachments = append(msg.Attachments, Attachment{
			Filename:    att.FileName,
			ContentType: att.ContentType,
			Size:        len(att.Content),
		})
	}

	return msg, nil
}

func parseAddressList(env *enmime.Envelope) []string {
	addrs, err := env.AddressList("To")
	if err != nil || len(addrs) == 0 {
		// Fall back to raw header
		if to := env.GetHeader("To"); to != "" {
			return []string{to}
		}
		return nil
	}
	result := make([]string, len(addrs))
	for i, addr := range addrs {
		result[i] = addr.String()
	}
	return result
}
