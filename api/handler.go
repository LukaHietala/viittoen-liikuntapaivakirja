package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/auth"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/db"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/services"
)

func NewHandler(s *db.Store, a *auth.Auth, m *services.MailService) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(a.Verifier())

	r.Route("/api", func(r chi.Router) {
		r.Mount("/collections", NewCollectionResource(s, a).Routes())
		r.Mount("/challenges", NewChallengeResource(s, a).Routes())
		r.Mount("/performances", NewPerformanceResource(s, a).Routes())
		r.Mount("/auth", NewAuthResource(s, a, m).Routes())
		r.Mount("/users", NewUserResource(s, a, m).Routes())
	})

	return r
}
