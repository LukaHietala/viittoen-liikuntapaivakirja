package main

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed templates/*.html
var templateFiles embed.FS

////go:embe static
//var staticFiles embed.FS

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

	tmpl, err := template.New("base").ParseFS(templateFiles, "templates/*.html")
	if err != nil {
		log.Fatalf("Unable to parse templates: %v\n", err)
	}

	app := &app{
		tmpl: tmpl,
	}

	r.Get("/", app.handleIndex)
	r.Get("/admin", app.handleAdmin)

	r.Route("/api", func(r chi.Router) {
		r.Route("/challenges", func(r chi.Router) {
			r.Get("/", GetChallenges)
			r.Post("/", CreateChallenge)
		})
	})

	return r
}


func GetChallenges(w http.ResponseWriter, r *http.Request) {
	challenges, err := getAllChallenges()
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	active := r.Form.Get("active")
	goal, err := strconv.Atoi(r.Form.Get("points"))

	if title == "" || err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	var isActive bool
	if active == "" {
		isActive = false
	} else {
		isActive = true
	}

	err = addChallenge(title, goal, isActive)
}

func (app *app) handleIndex(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	challenges, err := getAllChallenges()
	if err != nil {
		return
	}
	var pot int
	pot, err = getChallengePot(challenges[0])
	if err != nil {
		return
	}

	err = app.tmpl.ExecuteTemplate(&buf, "index.html", map[string]any{
		"Challenge": challenges[0],
		"Points":    pot,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}

	buf.WriteTo(w)
}

func (app *app) handleAdmin(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	challenges, err := getAllChallenges()
	if err != nil {
		return
	}

	err = app.tmpl.ExecuteTemplate(&buf, "challenges.html", map[string]any{
		"Challenges": challenges,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}

	buf.WriteTo(w)
}
