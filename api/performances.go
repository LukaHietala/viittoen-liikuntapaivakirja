package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/LukaHietala/viittoen-liikuntapaivakirja/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type performancesResource struct{}

func (rs performancesResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(SessionOnly)

	r.Post("/", rs.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(rs.PerformanceCtx)
		r.Delete("/", rs.Delete)
	})

	return r
}

func (rs performancesResource) PerformanceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var performance *db.Performance
		var err error

		rCtx := r.Context()

		performanceIDStr := chi.URLParam(r, "id")
		if performanceIDStr == "" {
			render.Render(w, r, ErrNotFound())
			return
		}

		performanceID, err := strconv.Atoi(performanceIDStr)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		performance, err = store.GetPerformanceByID(rCtx, performanceID)
		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), "performance", performance)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs performancesResource) Create(w http.ResponseWriter, r *http.Request) {
	var req db.Performance

	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := store.AddPerformance(r.Context(), &req); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	performance := &req
	render.Status(r, 201)
	render.JSON(w, r, performance)
}

func (rs performancesResource) Delete(w http.ResponseWriter, r *http.Request) {
	performance := r.Context().Value("performance").(*db.Performance)

	if err := store.DeletePerformance(r.Context(), performance.ID); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, 200)
	render.JSON(w, r, performance)
}
