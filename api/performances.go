package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/auth"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/db"
)

type performanceResource struct {
	store *db.Store
	auth  *auth.Auth
}

func NewPerformanceResource(s *db.Store, a *auth.Auth) *performanceResource {
	return &performanceResource{store: s, auth: a}
}

func (rs *performanceResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(rs.auth.Authenticator())

	r.Post("/", rs.Create)
	r.Route("/{performanceID}", func(r chi.Router) {
		r.Use(rs.PerformanceCtx)
		r.Delete("/", rs.Delete)
	})

	r.Get("/challenge/{challengeID}", rs.ListChallengePerformances)
	r.Get("/own", rs.ListOwnPerformances)
	// TODO: per user if needed

	return r
}

type PerformanceRequest struct {
	Points      int `json:"points"`
	ChallengeID int `json:"challengeId"`
}

func (p *PerformanceRequest) Bind(r *http.Request) error {
	if p.Points <= 0 {
		return errors.New("suorituksella tulee olla vähintään 1 piste")
	}
	if p.ChallengeID == 0 {
		return errors.New("haaste ei ole pätevä")
	}
	return nil
}

type PerformanceResponse struct {
	*db.Performance
}

func NewPerformanceResponse(performance *db.Performance) *PerformanceResponse {
	return &PerformanceResponse{Performance: performance}
}

func (p *PerformanceResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewPerformanceListResponse(performances []*db.Performance) []render.Renderer {
	list := []render.Renderer{}
	for _, performance := range performances {
		list = append(list, NewPerformanceResponse(performance))
	}
	return list
}

func (rs *performanceResource) PerformanceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		performanceIDStr := chi.URLParam(r, "performanceID")
		if performanceIDStr == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("performance id is missing")))
			return
		}

		performanceID, err := strconv.Atoi(performanceIDStr)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		performance, err := rs.store.FindPerformanceByID(r.Context(), performanceID)
		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), performanceContextKey, performance)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *performanceResource) ListOwnPerformances(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no valid user id")))
		return
	}

	performances, err := rs.store.ListUserPerformances(r.Context(), userID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.RenderList(w, r, NewPerformanceListResponse(performances))
}

func (rs *performanceResource) ListChallengePerformances(w http.ResponseWriter, r *http.Request) {
	challengeIDStr := chi.URLParam(r, "challengeID")
	if challengeIDStr == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("challenge id is missing")))
		return
	}

	challengeID, err := strconv.Atoi(challengeIDStr)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	challenge, err := rs.store.FindChallengeByID(r.Context(), challengeID)
	if err != nil {
		render.Render(w, r, ErrNotFound())
		return
	}

	performances, err := rs.store.ListChallengePerformances(r.Context(), challenge.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.RenderList(w, r, NewPerformanceListResponse(performances))
}

func (rs *performanceResource) Create(w http.ResponseWriter, r *http.Request) {
	data := &PerformanceRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userID, ok := getUserID(r)
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no valid user id")))
		return
	}

	newPerformance := &db.Performance{
		Points:      data.Points,
		ChallengeID: data.ChallengeID,
		UserID:      userID,
	}

	if err := rs.store.AddPerformance(r.Context(), newPerformance); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewPerformanceResponse(newPerformance))
}

func (rs *performanceResource) Delete(w http.ResponseWriter, r *http.Request) {
	performance, ok := r.Context().Value(performanceContextKey).(*db.Performance)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

	// Only admin or owner can delete a performance
	userID, ok := getUserID(r)
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no valid user id")))
		return
	}

	user, err := rs.store.FindUserByID(r.Context(), userID)
	if err != nil {
		render.Render(w, r, ErrUnauthorized())
		return
	}

	if !user.IsAdmin && performance.UserID != user.ID {
		render.Render(w, r, ErrUnauthorized())
		return
	}

	err = rs.store.DeletePerformanceByID(r.Context(), performance.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewPerformanceResponse(performance))
}
