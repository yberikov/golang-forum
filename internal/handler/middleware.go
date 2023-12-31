package handler

import (
	"context"
	"log"
	"net/http"

	"forum/internal/models"
)

type CtxKey string

const ctxKey CtxKey = "data"

func (h *Handler) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")

		var user models.User
		if err == http.ErrNoCookie {
			user = models.User{}
		} else if err == nil {
			user, err = h.services.Auth.GetUserByToken(cookie.Value)
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
		} else {
			h.ErrorPage(w, http.StatusInternalServerError, err)
		}

		ctx := context.WithValue(r.Context(), ctxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) isAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKey).(models.User)
		if user == (models.User{}) {
			http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}
