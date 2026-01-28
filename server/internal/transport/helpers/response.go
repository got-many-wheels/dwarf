package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	coreerror "github.com/got-many-wheels/dwarf/server/internal/core/error"
	"github.com/got-many-wheels/dwarf/server/internal/transport/middleware/logger"
	httpstatus "github.com/got-many-wheels/dwarf/server/internal/transport/status"
)

// WriteJSON writes V as the json value
func WriteJSON(ctx context.Context, status int, w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger := logger.FromContext(ctx)
		logger.Error(fmt.Sprintf("error encoding json: %v", err))
	}
}

// ErrorResponse represents a standard API error response
type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func WriteError(ctx context.Context, w http.ResponseWriter, log *slog.Logger, err error, msg string) {
	code := httpstatus.Status(err)

	if code >= http.StatusInternalServerError {
		log.ErrorContext(ctx, msg, "error", err)
	} else {
		log.WarnContext(ctx, msg, "error", err)
	}

	resp := ErrorResponse{
		Errors: make(map[string]string),
	}

	if derr, ok := coreerror.As(err); ok {
		switch {
		case len(derr.Fields) > 0:
			resp.Errors = derr.Fields
		case derr.Message != "":
			resp.Errors["root"] = derr.Message
		default:
			resp.Errors["root"] = http.StatusText(code)
		}
	} else {
		resp.Errors["root"] = http.StatusText(code)
	}
	WriteJSON(ctx, code, w, resp)
}
