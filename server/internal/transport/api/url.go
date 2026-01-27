package apiurl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/got-many-wheels/dwarf/server/internal/core"
	"github.com/got-many-wheels/dwarf/server/internal/service/url"
	"github.com/got-many-wheels/dwarf/server/internal/transport/middleware/logger"
)

func Register(mux *http.ServeMux, s *url.Service) {
	mux.HandleFunc("GET /{code}", get(s))
	mux.HandleFunc("POST /url", post(s))
	mux.HandleFunc("DELETE /url/{code}", delete(s))
}

func get(s *url.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logger.FromContext(ctx)
		code := r.PathValue("code")
		u, err := s.Get(context.Background(), code)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, u.Long, http.StatusSeeOther)
	}
}

func post(s *url.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logger.FromContext(ctx)
		var doc core.URL
		err := json.NewDecoder(r.Body).Decode(&doc)
		if err != nil {
			logger.Error(fmt.Sprintf("could not decode request body: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		doc.CreatedAt = time.Now().UTC()
		docs := []core.URL{doc}
		err = s.InsertBatch(context.Background(), docs)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(docs[0].String()))
	}
}

func delete(s *url.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logger.FromContext(ctx)
		code := r.PathValue("code")
		err := s.Delete(context.Background(), code)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
