package mux

import (
	"log/slog"
	"net/http"
	"os"

	services "github.com/got-many-wheels/dwarf/server/internal/service"
	apiurl "github.com/got-many-wheels/dwarf/server/internal/transport/api"
	"github.com/got-many-wheels/dwarf/server/internal/transport/middleware/logger"
	slogctx "github.com/veqryn/slog-context"
)

type Mux struct {
	*http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

func New(services services.Services) *Mux {
	slogBaseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})
	slogWithCtx := slogctx.NewHandler(slogBaseHandler, nil)
	l := slog.New(slogWithCtx)

	mux := &Mux{
		ServeMux: http.NewServeMux(),
	}
	mux.UseMiddleware(logger.Middleware(l))
	apiurl.Register(mux.ServeMux, services.URL)
	return mux
}

func (m *Mux) UseMiddleware(next func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, next)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler = m.ServeMux
	for _, next := range m.middlewares {
		h = next(h)
	}
	h.ServeHTTP(w, r)
}
