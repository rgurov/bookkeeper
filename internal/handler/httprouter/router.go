package httprouter

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rgurov/bookkeeper/internal/handler/httprouter/auth"
	"github.com/rgurov/bookkeeper/pkg/jwt"
	"github.com/rgurov/bookkeeper/pkg/logging"
	"golang.org/x/exp/slog"
)

func New(
	logger *slog.Logger,
	authService AuthService,
	jwtSecret string,
) http.Handler {
	jwt := jwt.NewJWT(jwtSecret)

	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		logging.Middleware(logger),
		middleware.Logger,
		middleware.Recoverer,
	)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", auth.NewRegister(jwt, authService))
		r.Post("/login", auth.NewLogin(jwt, authService))
	})

	return router
}
