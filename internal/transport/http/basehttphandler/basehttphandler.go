package basehttphandler

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	Logger        *slog.Logger
	CancelTimeout time.Duration
}

func (h *Handler) JSON(w http.ResponseWriter, status int, d any) {
	j, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write(j)
}

func Validate[T any](r *http.Request) (any, error) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err.Error() == "EOF" {
			return nil, errors.New("request body is empty")
		}
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}
	return req, nil
}
