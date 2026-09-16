package api

import (
	"encoding/json"
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
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type ErrResponse struct {
	Message string `json:"message"`
}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if req.Name == "" || req.Password == "" {
		res := ErrResponse{
			Message: "Nimi ja salasana ovat pakollisia",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(res)
		return
	}

	id, err := store.VerifyUser(req.Name, req.Password)
	if err != nil {
		res := ErrResponse{
			Message: "Nimi tai salasana on väärin",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(res)
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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		res := ErrResponse{
			Message: "Sähköposti ei ole oikeassa muodossa",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
		return
	}

	token := MakeResetToken(req.Email)
	err = store.StartResetPassword(req.Email, token, ms)
	if err != nil {
		res := ErrResponse{
			Message: "Käyttäjää ei löytynyt",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(res)
		return
	}
	w.WriteHeader(200)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		res := ErrResponse{
			Message: "Uusi salasana on pakollinen",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
		return
	}

	token, err := tokenAuth.Decode(req.Token)
	if err != nil {
		res := ErrResponse{
			Message: "Virheellinen tai vanhentunut nollaus linkki",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
		return
	}
	err = store.FinishResetPassword(token, req.NewPassword)
	if err != nil {
		res := ErrResponse{
			Message: "Jotain meni pieleen",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(200)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := store.ListUsers()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := store.GetSelf(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req db.User
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	if req.Name == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = store.UpdateUser(ctx, &req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	ctx := r.Context()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = store.DeleteUser(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func GetChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	challenge, err := store.GetChallengeByID(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}

func CreateUser(w http.ResponseWriter, r *http.Request, ms *services.MailService) {
	var req db.User

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	if req.Name == "" || req.Email == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = store.AddUser(&req, ms)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func GetChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challenges, err := store.ListChallenges(ctx, false)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

func GetActiveChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challenges, err := store.ListChallenges(ctx, true)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req db.Challenge

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	// TODO:
	if req.Title == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = store.AddChallenge(ctx, &req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func UpdateChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req db.Challenge

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	// TODO:
	if req.Title == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = store.UpdateChallenge(ctx, &req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func DeleteChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = store.DeleteChallenge(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func CreatePerformance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req db.Performance

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	err = store.AddPerformance(ctx, &req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func DeletePerformance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = store.DeletePerformance(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
