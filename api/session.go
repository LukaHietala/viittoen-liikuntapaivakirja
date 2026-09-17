package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type sessionResource struct{}

func (rs sessionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/self", rs.Self)
	r.Get("/valid", rs.Valid)

	return r
}

func (rs sessionResource) Self(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := store.GetSelf(ctx)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, user)
}

func (rs sessionResource) Valid(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]any)
	res["ok"] = true
	res["admin"] = false
	token, claims, _ := jwtauth.FromContext(r.Context())
	if token == nil || jwt.Validate(token) != nil {
		res["ok"] = false
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		res["ok"] = false

	}
	user, err := store.GetSession(int(userIDFloat))
	if err != nil || user == nil {
		res["ok"] = false
	}

	if user != nil {
		res["admin"] = user.IsAdmin
		res["name"] = user.Name
		res["id"] = user.ID
	}
	render.JSON(w, r, res)
}
