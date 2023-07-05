package logging

import (
	"net/http"

	"golang.org/x/exp/slog"
)

func Middleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(ContextWithLogger(r.Context(), logger))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
