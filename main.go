package main

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"encoding/json"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var db *sql.DB

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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")		
	})
	r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/challenges.html")		
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/challenges", func(r chi.Router) {
			r.Route("/active", func(r chi.Router) {
				r.Get("/", GetActiveChallenge)
			})
			r.Get("/", GetChallenges)
			r.Post("/", CreateChallenge)
			r.Route("/{id}", func (r chi.Router) {
				r.Delete("/{id}", DeleteChallenge)
			})
		})
		r.Route("/performances", func(r chi.Router) {
			r.Post("/", CreatePerformance)
		})
		r.Route("/users", func (r chi.Router)  {
			r.Get("/", GetUsers)
		})
	})

	return r
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
