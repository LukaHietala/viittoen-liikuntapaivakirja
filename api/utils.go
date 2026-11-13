package api

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func getUserID(r *http.Request) (int, bool) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	if id, ok := claims["user_id"].(float64); ok {
		return int(id), true
	}
	return 0, false
}
