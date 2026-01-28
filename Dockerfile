FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o smtp-gotify ./cmd/smtp-gotify

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/smtp-gotify /smtp-gotify

EXPOSE 2525 8080

ENTRYPOINT ["/smtp-gotify"]
