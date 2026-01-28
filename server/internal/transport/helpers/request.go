package helpers

import (
	"encoding/json"
	"net/http"

	httpstatus "github.com/got-many-wheels/dwarf/server/internal/transport/status"
)

func DecodeJSON[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, httpstatus.WithStatus(err, http.StatusBadRequest)
	}
	return v, nil
}
