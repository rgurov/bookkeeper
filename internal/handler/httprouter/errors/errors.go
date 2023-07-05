package errors

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"

	serviceErrors "github.com/rgurov/bookkeeper/internal/service/errors"
)

type httpError struct {
	StatusCode int    `json:"status_code"`
	Messages   string `json:"messages"`
}

func HttpError(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.WriteHeader(code)
	render.JSON(w, r, httpError{
		StatusCode: code,
		Messages:   err.Error(),
	})
}

func Abort(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, serviceErrors.ErrDataBusy):
		HttpError(w, r, http.StatusForbidden, err)
	case errors.Is(err, serviceErrors.ErrNotFound):
		HttpError(w, r, http.StatusNotFound, err)
	default:
		HttpError(w, r, http.StatusInternalServerError, err)
	}
}

func AbortWithBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	HttpError(w, r, http.StatusBadRequest, err)
}
