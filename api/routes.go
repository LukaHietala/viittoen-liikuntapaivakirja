package api

import (
	"errors"
	"io/fs"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/LukaHietala/viittoen-liikuntapaivakirja/db"
	"github.com/LukaHietala/viittoen-liikuntapaivakirja/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password"`
	Token       string `json:"token"`
}

var store *db.Store

func Serve(contentFS fs.FS, ms *services.MailService, s *db.Store) http.Handler {
	// :D
	store = s

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(jwtauth.Verifier(tokenAuth))

	staticFS, _ := fs.Sub(contentFS, "web/static")
	templateFS, _ := fs.Sub(contentFS, "web/templates")
	FileServer(r, "/static", staticFS)

	// Valid session required
	r.Group(func(r chi.Router) {
		r.Use(UnloggedInRedirector)

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, templateFS, "profile.html")
		})

		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			ResetJWTCookies(w)
			http.Redirect(w, r, "/", 303)
		})

		r.Get("/challenges", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, templateFS, "challenges.html")
		})
	})

	// No access with valid session
	r.Group(func(r chi.Router) {
		r.Use(LoggedInRedirector)
		r.Route("/login", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFileFS(w, r, templateFS, "login.html")
			})
			r.Post("/", Login)
		})

		r.Route("/forgot-password", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFileFS(w, r, templateFS, "forgot-password.html")
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				ForgotPassword(w, r, ms)
			})
		})

		r.Route("/reset-password", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFileFS(w, r, templateFS, "reset-password.html")
			})
			r.Post("/", ResetPassword)
		})
	})

	// Admin only pages
	r.Group(func(r chi.Router) {
		r.Use(AdminOnly)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, templateFS, "admin.html")
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, templateFS, "index.html")
	})

	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(SessionOnly)

			r.Route("/challenges", func(r chi.Router) {
				r.Get("/active", GetActiveChallenges)
				r.Get("/", GetChallenges)
				r.Get("/{id}", GetChallenge)
				r.Post("/", CreateChallenge)
				r.Patch("/", UpdateChallenge)
				r.Delete("/{id}", DeleteChallenge)
			})

			r.Route("/performances", func(r chi.Router) {
				r.Post("/", CreatePerformance)
				r.Delete("/{id}", DeletePerformance)
			})
			r.Route("/users", func(r chi.Router) {
				r.Use(AdminOnly)
				r.Get("/", GetUsers)
				r.Patch("/", UpdateUser)
				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					CreateUser(w, r, ms)
				})
				r.Delete("/{id}", DeleteUser)
			})
			r.Route("/session", func(r chi.Router) {
				r.Get("/", ValidSession)
			})
			r.Route("/self", func(r chi.Router) {
				r.Get("/", GetSelf)
			})
		})
	})

	return r
}

func FileServer(r chi.Router, path string, rootFS fs.FS) {
	root := http.FS(rootFS)
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func ValidSession(w http.ResponseWriter, r *http.Request) {
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

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	if req.Name == "" || req.Password == "" {
		render.Render(w, r, ErrInvalidRequest("Nimi ja salasana ovat pakollisia", errors.New("no name or password")))
		return
	}

	id, err := store.VerifyUser(req.Name, req.Password)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest("Nimi tai salasana on väärin", err))
		return
	}

	token := MakeSessionToken(id)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt",
		Value: token,
	})

	w.WriteHeader(200)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request, ms *services.MailService) {
	var req ForgotPasswordRequest
	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest("Sähköposti ei ole oikeassa muodossa", err))
		return
	}

	token := MakeResetToken(req.Email)
	err = store.StartResetPassword(req.Email, token, ms)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest("Käyttäjää ei löytynyt", err))
	}
	w.WriteHeader(200)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		// TODO: Erota
		render.Render(w, r, ErrInvalidRequest("Linkki on vanhentunut tai salasana on tyhjä", errors.New("no token or new password")))
		return
	}

	token, err := tokenAuth.Decode(req.Token)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest("Linkki on vanhentunut", err))
		return
	}
	err = store.FinishResetPassword(token, req.NewPassword)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	w.WriteHeader(200)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := store.ListUsers()
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, users)
}

func GetSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := store.GetSelf(ctx)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req db.User
	ctx := r.Context()
	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	if req.Name == "" {
		render.Render(w, r, ErrInvalidRequest("Nimi puuttuu", errors.New("name is missing")))
		return
	}

	err := store.UpdateUser(ctx, &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
	}
	w.WriteHeader(200)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	ctx := r.Context()
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	err = store.DeleteUser(ctx, id)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	w.WriteHeader(200)
}

func GetChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	challenge, err := store.GetChallengeByID(ctx, id)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, challenge)
}

func CreateUser(w http.ResponseWriter, r *http.Request, ms *services.MailService) {
	var req db.User
	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	if req.Name == "" || req.Email == "" {
		render.Render(w, r, ErrInvalidRequest("Nimi ja sähköposti ovat pakollisia", errors.New("no name or password")))
		return
	}

	err := store.AddUser(&req, ms)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return

	}
}

func GetChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challenges, err := store.ListChallenges(ctx, false)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, challenges)
}

func GetActiveChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challenges, err := store.ListChallenges(ctx, true)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, challenges)
}

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req db.Challenge

	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	// TODO:
	if req.Title == "" {
		render.Render(w, r, ErrInvalidRequest("Nimi on pakollinen", errors.New("no title")))
		return
	}

	err := store.AddChallenge(ctx, &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
}

func UpdateChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req db.Challenge

	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	// TODO:
	if req.Title == "" {
		render.Render(w, r, ErrInvalidRequest("Nimi on pakollinen", errors.New("title is missing")))
		return
	}

	err := store.UpdateChallenge(ctx, &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
}

func DeleteChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	err = store.DeleteChallenge(ctx, id)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
}

func CreatePerformance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req db.Performance

	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	err := store.AddPerformance(ctx, &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
}

func DeletePerformance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	err = store.DeletePerformance(ctx, id)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
}

type ErrResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(msg string, err error) render.Renderer {
	log.Println("invalid request error:", err)
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      msg,
	}
}

func ErrInternal(err error) render.Renderer {
	log.Println("internal error:", err)
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorText:      http.StatusText(http.StatusInternalServerError),
	}
}
