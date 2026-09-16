package api

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func ResetJWTCookies(w http.ResponseWriter) {
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

func LoggedInRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := jwtauth.FromContext(r.Context())

		if token != nil && jwt.Validate(token) == nil {
			http.Redirect(w, r, "/profile", 302)
		}

		next.ServeHTTP(w, r)
	})
}

func UnloggedInRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if token == nil || err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}

		if err := jwt.Validate(token); err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, _ := jwtauth.FromContext(r.Context())

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		user, err := store.GetUserByID(int(userIDFloat))

		if err != nil {
			ResetJWTCookies(w)
			http.Error(w, http.StatusText(403), 403)
			return
		}

		if !user.IsAdmin {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SessionOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, _ := jwtauth.FromContext(r.Context())

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		user, err := store.GetUserByID(int(userIDFloat))

		if err != nil {
			ResetJWTCookies(w)
			http.Error(w, http.StatusText(403), 403)
			return
		}

		if user == nil {
			ResetJWTCookies(w)
			http.Error(w, http.StatusText(403), 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}
