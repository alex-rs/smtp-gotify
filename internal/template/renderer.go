package template

import (
	"bytes"
	"text/template"

	"github.com/alex/smtp-gotify/internal/mail"
)

type TemplateData struct {
	From    string
	To      string
	Subject string
	Body    string
}

type Renderer struct {
	titleTpl   *template.Template
	messageTpl *template.Template
}

func NewRenderer(titleTemplate, messageTemplate string) (*Renderer, error) {
	titleTpl, err := template.New("title").Parse(titleTemplate)
	if err != nil {
		return nil, err
	}

	messageTpl, err := template.New("message").Parse(messageTemplate)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		titleTpl:   titleTpl,
		messageTpl: messageTpl,
	}, nil
}

func (r *Renderer) Render(msg *mail.Message) (title, message string, err error) {
	data := TemplateData{
		From:    msg.From,
		To:      joinAddresses(msg.To),
		Subject: msg.Subject,
		Body:    msg.Body,
	}

	var titleBuf bytes.Buffer
	if err := r.titleTpl.Execute(&titleBuf, data); err != nil {
		return "", "", err
	}

	var msgBuf bytes.Buffer
	if err := r.messageTpl.Execute(&msgBuf, data); err != nil {
		return "", "", err
	}

	return titleBuf.String(), msgBuf.String(), nil
}

func joinAddresses(addrs []string) string {
	if len(addrs) == 0 {
		return ""
	}
	result := addrs[0]
	for i := 1; i < len(addrs); i++ {
		result += ", " + addrs[i]
	}
	return result
}
