package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"net/url"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

var db *sql.DB
var tokenAuth *jwtauth.JWTAuth

const Secret = "suuri_salaisuus"

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
}

func makeToken(userId int) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": userId})
	return tokenString
}

func main() {
	server := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: serve(),
	}

	var err error
	db, err = connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}

func serve() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Valid session required
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Use(UnloggedInRedirector)

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("Hei siellä %v", claims["user_id"])))
		})

		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				HttpOnly: true,
				MaxAge:   -1,
				SameSite: http.SameSiteLaxMode,
				// Uncomment below for HTTPS:
				// Secure: true,
				Name:  "jwt",
				Value: "",
			})

			http.Redirect(w, r, "/", 303)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Use(LoggedInRedirector)
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/login.html")
		})
		
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
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

			id, err := verifyUser(name, password)
			if err != nil {
				params := url.Values{}
				params.Add("err", "invalid")
				finalURL := "/login" + "?" + params.Encode()
				http.Redirect(w, r, finalURL, 303)
				return
			}

			token := makeToken(id)

			http.SetCookie(w, &http.Cookie{
				HttpOnly: true,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				SameSite: http.SameSiteLaxMode,
				// Uncomment below for HTTPS:
				// Secure: true,
				Name:  "jwt",
				Value: token,
			})

			http.Redirect(w, r, "/profile", 303)
		})
	})

	r.Group(func (r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(AdminOnly)
		
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/challenges.html")
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	

	r.Route("/api", func(r chi.Router) {
		r.Route("/challenges", func(r chi.Router) {
			r.Route("/active", func(r chi.Router) {
				r.Get("/", GetActiveChallenge)
			})
			r.Get("/", GetChallenges)
			r.Post("/", CreateChallenge)
			r.Route("/{id}", func(r chi.Router) {
				r.Delete("/{id}", DeleteChallenge)
			})
		})
		r.Route("/performances", func(r chi.Router) {
			r.Post("/", CreatePerformance)
		})
		r.Route("/users", func(r chi.Router) {
			r.Get("/", GetUsers)
		})
	})

	return r
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
		token, _, _ := jwtauth.FromContext(r.Context())

		if token == nil || jwt.Validate(token) != nil {
			http.Redirect(w, r, "/login", 302)
		}

		next.ServeHTTP(w, r)
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, _ := jwtauth.FromContext(r.Context())
		
		if token == nil || jwt.Validate(token) != nil {
			http.Redirect(w, r, "/", 302)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Redirect(w, r, "/", 302)
			return
		}

		user, err := getUser(int(userIDFloat))

		if err != nil {
			http.Redirect(w, r, "/", 302)
			return
		}

		if !user.IsAdmin {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}



func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := getAllUsers()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetActiveChallenge(w http.ResponseWriter, r *http.Request) {
	c, err := getActiveChallenge()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func GetChallenges(w http.ResponseWriter, r *http.Request) {
	challenges, err := getAllChallenges()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	var req Challenge

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	if req.Title == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = addChallenge(req.Title, req.GoalPoints, req.IsActive)
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

	err = deleteChallenge(id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func CreatePerformance(w http.ResponseWriter, r *http.Request) {
	var req Performance

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	err = addPerformance(req.Points, req.UserID, req.ChallengeID)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
