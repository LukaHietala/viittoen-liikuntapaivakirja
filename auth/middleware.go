package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (a *Auth) RejectAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil {
			next.ServeHTTP(w, r)
			return
		}

		if err := jwt.Validate(token); err == nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Auth) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil {
			next.ServeHTTP(w, r)
			return
		}

		if err := jwt.Validate(token); err == nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func resetJWTCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt",
		Value: "",
	})
}
