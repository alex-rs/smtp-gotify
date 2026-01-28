# smtp-gotify

[![CI](https://github.com/alex-rs/smtp-gotify/actions/workflows/ci.yml/badge.svg)](https://github.com/alex-rs/smtp-gotify/actions/workflows/ci.yml)

A lightweight SMTP server that forwards incoming emails as Gotify push notifications.

## Use Case

Many applications and services can send email notifications but lack support for modern push notification systems. smtp-gotify bridges this gap by acting as an SMTP server that converts emails into Gotify notifications, enabling you to receive push alerts from legacy systems, IoT devices, cron jobs, or any software that supports email alerts.

## Features

- SMTP server with ESMTP support (RFC 5321)
- MIME email parsing (plain text and HTML)
- Multiple Gotify tokens (broadcast to multiple devices/apps)
- Customizable notification templates
- Optional markdown rendering
- Health check endpoint
- Structured JSON logging
- Minimal Docker image (~10MB)

## Configuration

All configuration is done via environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GOTIFY_URL` | Yes | - | Gotify server URL |
| `GOTIFY_TOKEN` | Yes | - | App token(s), comma-separated for multiple |
| `GOTIFY_PRIORITY` | No | `5` | Message priority (0-10) |
| `GOTIFY_MARKDOWN` | No | `false` | Enable markdown rendering |
| `GOTIFY_TITLE_TEMPLATE` | No | `{{.Subject}}` | Notification title template |
| `GOTIFY_MESSAGE_TEMPLATE` | No | See below | Notification body template |
| `SMTP_LISTEN` | No | `:2525` | SMTP listen address |
| `SMTP_DOMAIN` | No | `localhost` | SMTP domain name |
| `SMTP_MAX_SIZE` | No | `10485760` | Max message size (bytes) |
| `HEALTH_ENABLED` | No | `true` | Enable health endpoint |
| `HEALTH_LISTEN` | No | `:8080` | Health endpoint address |
| `LOG_LEVEL` | No | `info` | Log level (debug/info/warn/error) |
| `LOG_FORMAT` | No | `json` | Log format (json/text) |

Default message template:
```
From: {{.From}}
To: {{.To}}
---
{{.Body}}
```

### Template Variables

- `{{.From}}` - Sender address
- `{{.To}}` - Recipient address(es)
- `{{.Subject}}` - Email subject
- `{{.Body}}` - Email body (plain text preferred, falls back to HTML)

## Quick Start

1. Download the compose file:
```bash
curl -O https://raw.githubusercontent.com/alex-rs/smtp-gotify/main/docker-compose.yml
```

2. Edit `docker-compose.yml` and set your Gotify URL and token:
```yaml
environment:
  GOTIFY_URL: "https://gotify.example.com"
  GOTIFY_TOKEN: "your-app-token"
```

3. Start the service:
```bash
docker compose up -d
```

For multiple Gotify tokens (sends to all):
```yaml
environment:
  GOTIFY_TOKEN: "token1,token2,token3"
```

## Testing

Send a test email using sendmail, swaks, or any SMTP client:

```bash
# Using sendmail
echo -e "Subject: Test Alert\n\nThis is a test message." | sendmail -S localhost:2525 test@example.com

# Using swaks
swaks --to test@example.com --from sender@example.com --server localhost:2525 --body "Test message"

# Using curl (health check)
curl http://localhost:8080/health
```

## Building

```bash
# Build binary
make build

# Run tests
make test

# Build Docker image
make docker
```

## License

MIT
