# Guides

## Folder structures

```bash
.
├── cmd             # Binary entry point.
├── docs            # Auto generated OpenAPI (Swagger) documents (don't touch this).
├── internal
│   ├── core        # Domain and business invariants to represents concepts.
│   ├── platform    # "Duct Tape" that wires the whole systems.
│   ├── service     # Application uses cases and orchestration.
│   ├── store       # Persistence Adapter that translates domain from storage
│   └── transport   # Where the http transport lives, duh.
└── utils
```

### Run the server

```bash
go mod download

# Required to generate sql queries to code/business logic
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go generate ./cmd

export DATABASE_URI=...
export PORT=:8080
go run ./cmd
```

### Generate swagger API endpoints documents

```bash
swagger -g cmd/main.go
```
