package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/LukaHietala/viittoen-liikuntapaivakirja/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type challengesResource struct{}

func (rs challengesResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(SessionOnly)

	r.Get("/", rs.List)
	r.Post("/", rs.Create)
	r.Get("/active", rs.GetActive)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(rs.ChallengeCtx)
		r.Get("/", rs.Get)
		r.Delete("/", rs.Delete)
		r.Patch("/", rs.Update)
	})

	return r
}

func (rs challengesResource) ChallengeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var challenge *db.Challenge
		var err error

		rCtx := r.Context()

		challengeIDStr := chi.URLParam(r, "id")
		if challengeIDStr == "" {
			render.Render(w, r, ErrNotFound())
		}

		challengeID, err := strconv.Atoi(challengeIDStr)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
		}

		challenge, err = store.GetChallengeByID(rCtx, challengeID)
		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), "challenge", challenge)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs challengesResource) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challenges, err := store.ListChallenges(ctx, false)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, challenges)
}

func (rs challengesResource) Create(w http.ResponseWriter, r *http.Request) {
	var challenge db.Challenge

	if err := render.Decode(r, &challenge); err != nil {
		render.Render(w, r, ErrInvalidRequest("invalid json payload", err))
		return
	}

	// TODO:
	if challenge.Title == "" {
		render.Render(w, r, ErrInvalidRequest("Nimi on pakollinen", errors.New("no title")))
		return
	}

	if err := store.AddChallenge(r.Context(), challenge); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, 201)
	render.JSON(w, r, challenge)
}

func (rs challengesResource) Get(w http.ResponseWriter, r *http.Request) {
	challenge := r.Context().Value("challenge").(*db.Challenge)
	render.JSON(w, r, challenge)
}

func (rs challengesResource) GetActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challenges, err := store.ListChallenges(ctx, true)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, challenges)
}

func (rs challengesResource) Update(w http.ResponseWriter, r *http.Request) {
	challenge := r.Context().Value("challenge").(*db.Challenge)

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

	challenge = &req

	if err := store.UpdateChallenge(r.Context(), challenge.ID, req); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
}

func (rs challengesResource) Delete(w http.ResponseWriter, r *http.Request) {
	challenge := r.Context().Value("challenge").(*db.Challenge)

	err := store.DeleteChallenge(r.Context(), challenge.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, 200)
	render.JSON(w, r, challenge)
}
