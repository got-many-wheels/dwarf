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

### Generate swagger API endpoints documents

```bash
swagger -g cmd/main.go
```
