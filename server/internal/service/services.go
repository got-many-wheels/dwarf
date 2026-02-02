package services

import (
	"github.com/got-many-wheels/dwarf/server/internal/service/url"
)

// Services bundles the use-case layer for consumption of transports
type Services struct {
	URL *url.Service
}

// Stores provides the persistence interfaces required to build Services
type Stores struct {
	URL url.Store
}

// Wires repositories into use-case services and connects cross-cutting dependencies
func New(stores Stores) Services {
	return Services{
		URL: url.New(stores.URL),
	}
}
