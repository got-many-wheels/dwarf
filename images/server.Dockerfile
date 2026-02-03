# Golang Server binary build step
FROM golang:1.24-bookworm AS builder

WORKDIR /usr/local/server

COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy the rest of the server files
COPY server .

# sqlc needs schema.sql during go generate, copy the whole migrations folder. (lol)
COPY migrations ../migrations
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN go generate ./...
RUN CGO_ENABLED=0 go build -v -o /usr/local/server/bin/app ./cmd

# Running step
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /usr/local/server/bin/app ./server

EXPOSE 8080
ENTRYPOINT ["/app/server"]
