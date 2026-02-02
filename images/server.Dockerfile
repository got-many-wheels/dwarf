# Golang Server binary build step
FROM golang:1.24-bookworm AS builder
WORKDIR /usr/local/server
COPY server .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/server/bin/app ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /usr/local/server/bin/app ./server

EXPOSE 8080
ENTRYPOINT ["/app/server"]
