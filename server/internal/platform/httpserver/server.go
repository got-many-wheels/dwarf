package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// http server timeouts
const (
	writetimeout = 10 * time.Second
	readtimeout  = 10 * time.Second
	idletimeout  = 1 * time.Minute
)

type Server struct {
	server *http.Server
}

func New(handler http.Handler, addr string) *Server {
	return &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			WriteTimeout: writetimeout,
			ReadTimeout:  readtimeout,
			IdleTimeout:  idletimeout,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("server is up and running on http://localhost%s", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				return
			}
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		log.Println("http server closed successfully")
	case err := <-errCh:
		return fmt.Errorf("httpserver.ListenAndServe error: %v", err)
	}
	return nil
}
