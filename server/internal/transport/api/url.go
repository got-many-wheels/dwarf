package apiurl

import (
	"context"
	"net/http"
	"time"

	coreurl "github.com/got-many-wheels/dwarf/server/internal/core/url"
	"github.com/got-many-wheels/dwarf/server/internal/service/url"
	"github.com/got-many-wheels/dwarf/server/internal/transport/helpers"
	"github.com/got-many-wheels/dwarf/server/internal/transport/middleware/logger"
)

func Register(mux *http.ServeMux, s *url.Service) {
	mux.HandleFunc("GET /{code}", get(s))
	mux.HandleFunc("POST /url", post(s))
	mux.HandleFunc("DELETE /url/{code}", delete(s))
}

// get godoc
// @Summary Redirect to original URL
// @Description Resolve a short code and redirect to the original long URL
// @Tags url
// @Param code path string true "Short URL code"
// @Success 303 "Redirect"
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /{code} [get]
func get(s *url.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logger.FromContext(ctx)
		code := r.PathValue("code")
		u, err := s.Get(context.Background(), code)
		if err != nil {
			helpers.WriteError(ctx, w, logger, err, "s.Get() error")
			return
		}
		http.Redirect(w, r, u.Long, http.StatusSeeOther)
	}
}

// post godoc
// @Summary Create a shortened URL
// @Description Create a new short URL entry
// @Tags url
// @Accept json
// @Produce json
// @Param url body coreurl.URL true "URL payload"
// @Success 201 {object} coreurl.URL
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 409 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /url [post]
func post(s *url.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logger.FromContext(ctx)
		doc, err := helpers.DecodeJSON[coreurl.URL](r)
		if err != nil {
			helpers.WriteError(ctx, w, logger, err, "helpers.DecodeJSON error")
			return
		}
		doc.CreatedAt = time.Now().UTC()
		docs := []coreurl.URL{doc}
		err = s.InsertBatch(context.Background(), docs)
		if err != nil {
			helpers.WriteError(ctx, w, logger, err, "s.InsertBatch() error")
			return
		}

		helpers.WriteJSON(ctx, http.StatusCreated, w, docs[0]) // index 0 because we only insert once
	}
}

// delete godoc
// @Summary Delete a shortened URL
// @Description Delete a short URL by code
// @Tags url
// @Param code path string true "Short URL code"
// @Success 204 "No Content"
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /url/{code} [delete]
func delete(s *url.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logger.FromContext(ctx)
		code := r.PathValue("code")
		err := s.Delete(context.Background(), code)
		if err != nil {
			helpers.WriteError(ctx, w, logger, err, "s.Delete() error")
			return
		}
		helpers.WriteJSON(ctx, http.StatusNoContent, w, nil)
	}
}
