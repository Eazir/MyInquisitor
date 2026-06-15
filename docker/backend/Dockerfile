ARG GO_VERSION=1.26

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app ./cmd/api

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
COPY --from=builder /src/internal/infrastructure/persistence/migrations /migrations

EXPOSE 8080

ENTRYPOINT ["/app"]
