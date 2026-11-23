package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (a *Auth) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, err := a.store.FindUserByID(r.Context(), int(userIDFloat))

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !user.IsAdmin {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
