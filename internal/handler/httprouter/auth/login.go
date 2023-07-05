package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/rgurov/bookkeeper/internal/handler/httprouter/errors"
	"github.com/rgurov/bookkeeper/pkg/jwt"
)

type LoginService interface {
	Login(ctx context.Context, login, password string) (int, error)
}

type loginDTO struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func NewLogin(jwt *jwt.JWT, login LoginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto loginDTO
		if err := render.DecodeJSON(r.Body, &dto); err != nil {
			errors.AbortWithBadRequest(w, r, err)
			return
		}

		id, err := login.Login(r.Context(), dto.Login, dto.Password)
		if err != nil {
			errors.Abort(w, r, err)
			return
		}

		token, err := jwt.Sign(
			map[string]interface{}{
				"sub": id,
			},
			time.Hour*24*365, // 1 year
		)
		if err != nil {
			errors.Abort(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, loginResponse{
			AccessToken: token,
		})
	}
}
