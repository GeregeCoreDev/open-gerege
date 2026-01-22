# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /app/bin/server ./cmd/server/main.go

# ---------- Runtime stage ----------
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/bin/server /app/server

USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/server"]
