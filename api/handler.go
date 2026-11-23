package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/auth"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/db"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/services"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/web"
)

func NewHandler(s *db.Store, a *auth.Auth, m *services.MailService) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(a.Verifier())

	r.Route("/api", func(r chi.Router) {
		r.Mount("/collections", NewCollectionResource(s, a).Routes())
		r.Mount("/challenges", NewChallengeResource(s, a, m).Routes())
		r.Mount("/performances", NewPerformanceResource(s, a).Routes())
		r.Mount("/auth", NewAuthResource(s, a, m).Routes())
		r.Mount("/users", NewUserResource(s, a, m).Routes())
	})

	dist, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(dist))

	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		if path == "" {
			path = "index.html"
		}

		if f, err := dist.Open(path); err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFileFS(w, r, dist, "index.html")
	}))
	return r
}
