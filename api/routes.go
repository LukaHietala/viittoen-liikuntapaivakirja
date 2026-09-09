package api

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"net/mail"
	"net/url"
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

func Serve(contentFS fs.FS, ms *services.MailService) http.Handler {
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
			// TODO: Rename html
			http.ServeFileFS(w, r, templateFS, "chanllenges.html")
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
				r.Get("/latest", GetLatestChallenge)
				r.Get("/", GetChallenges)
				r.Post("/", CreateChallenge)
				r.Patch("/", UpdateChallenge)
				r.Get("/{id}", GetChallenge)
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
	user, err := db.GetSession(int(userIDFloat))
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
	r.ParseForm()
	name := r.PostForm.Get("name")
	password := r.PostForm.Get("password")

	if name == "" || password == "" {
		params := url.Values{}
		params.Add("err", "missing")
		finalURL := "/login" + "?" + params.Encode()
		http.Redirect(w, r, finalURL, 303)
		return
	}

	id, err := db.VerifyUser(name, password)
	if err != nil {
		params := url.Values{}
		params.Add("err", "invalid")
		finalURL := "/login" + "?" + params.Encode()
		http.Redirect(w, r, finalURL, 303)
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

	http.Redirect(w, r, "/", 303)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request, ms *services.MailService) {
	r.ParseForm()
	email := r.PostForm.Get("email")

	_, err := mail.ParseAddress(email)
	if err != nil {
		params := url.Values{}
		params.Add("err", "invalid")
		finalURL := "/forgot-password" + "?" + params.Encode()
		http.Redirect(w, r, finalURL, 303)
		return
	}

	token := MakeResetToken(email)
	err = db.StartResetPassword(email, token, ms)
	if err != nil {
		params := url.Values{}
		params.Add("err", "failed")
		finalURL := "/forgot-password" + "?" + params.Encode()
		http.Redirect(w, r, finalURL, 303)
		return
	}

	params := url.Values{}
	params.Add("state", "success")
	finalURL := "/forgot-password" + "?" + params.Encode()
	http.Redirect(w, r, finalURL, 303)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	password := r.PostForm.Get("password")
	tokenStr := r.PostForm.Get("token")

	if tokenStr == "" || password == "" {
		params := url.Values{}
		params.Add("token", tokenStr)
		params.Add("err", "invalid")
		finalURL := "/reset-password" + "?" + params.Encode()
		http.Redirect(w, r, finalURL, 303)
		return
	}

	token, err := tokenAuth.Decode(tokenStr)
	if err != nil {
		// TODO:
		log.Println(err)
		http.Redirect(w, r, "/", 303)
		return
	}
	err = db.FinishResetPassword(token, password)
	if err != nil {
		params := url.Values{}
		params.Add("token", tokenStr)
		params.Add("err", "failed")
		finalURL := "/reset-password" + "?" + params.Encode()
		http.Redirect(w, r, finalURL, 303)
		return
	}

	http.Redirect(w, r, "/login", 303)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetAllUsers()
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
	user, err := db.GetSelf(ctx)
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

	err = db.UpdateUser(ctx, &req)
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

	err = db.DeleteUser(ctx, id)
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
	challenge, err := db.GetChallenge(ctx, id)
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

	if req.Name == "" || req.PasswordPlain == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = db.AddUser(&req, ms)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func GetChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challenges, err := db.GetAllChallenges(ctx, false)
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
	challenges, err := db.GetAllChallenges(ctx, true)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
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

	err = db.AddChallenge(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func UpdateChallenge(w http.ResponseWriter, r *http.Request) {
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

	err = db.UpdateChallenge(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func DeleteChallenge(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = db.DeleteChallenge(id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func GetLatestChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challenge, err := db.GetLatestChallenge(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
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

	err = db.AddPerformance(ctx, &req)
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

	err = db.DeletePerformance(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

