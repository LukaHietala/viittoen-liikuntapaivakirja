package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"html/template"
	"bytes"
	"embed"
	"strconv"
	"slices"
	

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Challenge struct {
	ID int `json:"id"`
	Title string `json:"title"`
	GoalPoints int `json:"goal_points"`
	IsActive bool `json:"is_active"`
	Performances []*Performance `json:"performances"`
}

func (c *Challenge) totalPoints() (points int) {
	for _, p := range c.Performances {
		points += p.Points
	}
	return
}

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Points int `json:"points"`
	PasswordHash string `json:"password_hash"`
	IsAdmin bool `json:"is_admin"`
	Performances []*Performance `json:"performances"`
}

type Performance struct {
	ID int `json:"id"`
	Points int `json:"points"`
	ChallengeID int `json:"challenge_id"`
}

var users = []*User{
	{ID: 1, Name: "Jaakko", Points: 12, PasswordHash: "1234", IsAdmin: true},
	{ID: 2, Name: "Tero", Points: 10, PasswordHash: "1224", IsAdmin: true},
	{ID: 3, Name: "Jorma", Points: 102, PasswordHash: "1334", IsAdmin: false},
}

var challenges = []*Challenge{
	{ID: 1, Title: "Haaste 1", GoalPoints: 12, IsActive: true},
	{ID: 2, Title: "Haaste 2", GoalPoints: 122, IsActive: false},
}

//go:embed templates/*.html
var templateFiles embed.FS

////go:embe static
//var staticFiles embed.FS

type app struct {
	tmpl *template.Template
}

func main() {
	server := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: serve(),
	}

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
			r.Post("/", CreateChallenge)
		})
	})

	return r
}

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	goal, err := strconv.Atoi(r.Form.Get("points"))

	if name == "" || err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	challenges = append(challenges,	&Challenge{ID: 15, Title: name, GoalPoints: goal, IsActive: true})
}

func (app *app) handleIndex(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	activeIndex := slices.IndexFunc(challenges, func(c *Challenge) bool {
		return c.IsActive
	})
	
	err := app.tmpl.ExecuteTemplate(&buf, "index.html", map[string]any{
		"Challenge": challenges[activeIndex],
		"Points": challenges[activeIndex].totalPoints(),
	})

    if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}

	buf.WriteTo(w)
}

func (app *app) handleAdmin(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	err := app.tmpl.ExecuteTemplate(&buf, "challenges.html", map[string]any{
		"Challenges": challenges,
		"Users": users,
	})
    if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}

	buf.WriteTo(w)
}