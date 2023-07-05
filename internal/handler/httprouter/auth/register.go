package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/rgurov/bookkeeper/internal/handler/httprouter/errors"
	"github.com/rgurov/bookkeeper/pkg/jwt"
)

type RegisterService interface {
	Register(ctx context.Context, login, password string) (int, error)
}

type registerDTO struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerResponse struct {
	AccessToken string `json:"access_token"`
}

func NewRegister(jwt *jwt.JWT, register RegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto registerDTO
		if err := render.DecodeJSON(r.Body, &dto); err != nil {
			errors.AbortWithBadRequest(w, r, err)
			return
		}
		id, err := register.Register(r.Context(), dto.Login, dto.Password)
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
		render.JSON(w, r, registerResponse{
			AccessToken: token,
		})
	}
}
