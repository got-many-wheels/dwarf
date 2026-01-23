package apiurl

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/got-many-wheels/dwarf/server/internal/core"
	"github.com/got-many-wheels/dwarf/server/internal/service/url"
)

func Register(mux *http.ServeMux, s *url.URLService) {
	mux.HandleFunc("GET /{code}", get(s))
	mux.HandleFunc("POST /url", post(s))
	mux.HandleFunc("DELETE /url/{code}", delete(s))
}

func get(s *url.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		u, err := s.Get(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, u.Long, http.StatusSeeOther)
	}
}

func post(s *url.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var doc core.URL
		err := json.NewDecoder(r.Body).Decode(&doc)
		if err != nil {
			log.Printf("could not decode request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		doc.CreatedAt = time.Now().UTC()
		err = s.InsertBatch(context.Background(), []core.URL{doc})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(doc.String()))
	}
}

func delete(s *url.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		err := s.Delete(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
