package mux

import (
	"net/http"

	services "github.com/got-many-wheels/dwarf/server/internal/service"
	apiurl "github.com/got-many-wheels/dwarf/server/internal/transport/api"
)

type Mux struct {
	*http.ServeMux
}

func New(services services.Services) *Mux {
	mux := &Mux{
		ServeMux: http.NewServeMux(),
	}
	apiurl.Register(mux.ServeMux, services.URL)
	return mux
}
