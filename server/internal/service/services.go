package services

import "github.com/got-many-wheels/dwarf/server/internal/service/url"

// Services bundles the use-case layer for consumption of transports
type Services struct {
	URL *url.URLService
}

// Stores provides the persistence interfaces required to build Services
type Stores struct {
	URL url.URLStore
}

// Wires repositories into use-case services and connects cross-cutting dependencies
func New(stores Stores) Services {
	return Services{
		URL: url.NewURLService(stores.URL),
	}
}
