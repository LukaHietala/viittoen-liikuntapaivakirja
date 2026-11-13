package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/auth"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/db"
)

type challengeResource struct {
	store *db.Store
	auth  *auth.Auth
}

func NewChallengeResource(s *db.Store, a *auth.Auth) *challengeResource {
	return &challengeResource{store: s, auth: a}
}

func (rs *challengeResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(rs.auth.Authenticator())

	r.Get("/", rs.List)
	r.Get("/active", rs.ListActive)
	r.Post("/", rs.Create)

	r.Route("/{challengeID}", func(r chi.Router) {
		r.Use(rs.ChallengeCtx)
		r.Get("/", rs.FindOne)
		r.Delete("/", rs.Delete)
		r.Put("/", rs.Update)
	})

	return r
}

type ChallengeRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Unit         string `json:"unit"`
	GoalPoints   int    `json:"goalPoints"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	CollectionID int    `json:"collectionId"`
}

func (c *ChallengeRequest) Bind(r *http.Request) error {
	c.Title = strings.TrimSpace(c.Title)
	c.Description = strings.TrimSpace(c.Description)
	c.Unit = strings.TrimSpace(c.Unit)

	if c.Title == "" {
		return errors.New("haasteen nimi on pakollinen")
	}
	if utf8.RuneCountInString(c.Title) > 50 {
		return errors.New("haasteen nimi on liian pitkä (max 50)")
	}

	if c.Description == "" {
		return errors.New("haasteen kuvaus on pakollinen")
	}
	if utf8.RuneCountInString(c.Description) > 2000 {
		return errors.New("haasteen kuvaus on liian pitkä (max 2000)")
	}

	if c.Unit == "" {
		return errors.New("haasteen suoritusyksikkö on pakollinen")
	}
	if utf8.RuneCountInString(c.Unit) > 20 {
		return errors.New("haasteen suoritusyksikkö on liian pitkä (max 20)")
	}

	if c.GoalPoints < 1 {
		return errors.New("haasteen tavoitepisteet tulee olla yli 1")
	}

	start, err := time.Parse(time.DateOnly, c.StartDate)
	if err != nil {
		return errors.New("aloituspäivä on väärässä formaatissa")
	}
	end, err := time.Parse(time.DateOnly, c.EndDate)
	if err != nil {
		return errors.New("lopetuspäivä on väärässä formaatissa")
	}
	if !end.After(start) {
		return errors.New("aloituspäivän tulee olla aina ennen lopetuspäivää")
	}

	return nil
}

type ChallengeResponse struct {
	*db.Challenge
	Performances []*db.Performance `json:"performances"`
	Pot          int               `json:"pot"`
}

func NewChallengeResponse(challenge *db.Challenge, performances []*db.Performance) *ChallengeResponse {
	return &ChallengeResponse{Challenge: challenge, Performances: performances}
}

func (c *ChallengeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	for _, p := range c.Performances {
		c.Pot += p.Points
	}
	return nil
}

func NewChallengeListResponse(challenges []*db.Challenge, performances []*db.Performance) []render.Renderer {
	list := []render.Renderer{}
	for _, c := range challenges {
		cp := make([]*db.Performance, 0)
		for _, p := range performances {
			if p.ChallengeID == c.ID {
				cp = append(cp, p)
			}
		}
		list = append(list, NewChallengeResponse(c, cp))
	}
	return list
}

func (rs *challengeResource) ChallengeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), challengeContextKey, challenge)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *challengeResource) List(w http.ResponseWriter, r *http.Request) {
	challenges, err := rs.store.ListChallenges(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	performances, err := rs.store.ListPerformances(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.RenderList(w, r, NewChallengeListResponse(challenges, performances))
}

func (rs *challengeResource) ListActive(w http.ResponseWriter, r *http.Request) {
	challenges, err := rs.store.ListActiveChallenges(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	performances, err := rs.store.ListPerformances(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.RenderList(w, r, NewChallengeListResponse(challenges, performances))
}

func (rs *challengeResource) Create(w http.ResponseWriter, r *http.Request) {
	data := &ChallengeRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userID, ok := getUserID(r)
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no valid user id")))
		return
	}

	newChallenge := &db.Challenge{
		Title:        data.Title,
		Description:  data.Description,
		Unit:         data.Unit,
		GoalPoints:   data.GoalPoints,
		StartDate:    data.StartDate,
		EndDate:      data.EndDate,
		UserID:       userID,
		CollectionID: data.CollectionID,
	}

	if err := rs.store.AddChallenge(r.Context(), newChallenge); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewChallengeResponse(newChallenge, nil))
}

func (rs *challengeResource) FindOne(w http.ResponseWriter, r *http.Request) {
	challenge, ok := r.Context().Value(challengeContextKey).(*db.Challenge)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

	performances, err := rs.store.ListChallengePerformances(r.Context(), challenge.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewChallengeResponse(challenge, performances))
}

func (rs *challengeResource) Update(w http.ResponseWriter, r *http.Request) {
	currentChallenge, ok := r.Context().Value(challengeContextKey).(*db.Challenge)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

	data := &ChallengeRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	currentChallenge.Title = data.Title
	currentChallenge.Description = data.Description
	currentChallenge.Unit = data.Unit
	currentChallenge.GoalPoints = data.GoalPoints
	currentChallenge.StartDate = data.StartDate
	currentChallenge.EndDate = data.EndDate

	// Only admin or owner can edit the challenge
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

	if !user.IsAdmin && currentChallenge.UserID != user.ID {
		render.Render(w, r, ErrUnauthorized())
		return
	}

	if err := rs.store.UpdateChallenge(r.Context(), currentChallenge); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewChallengeResponse(currentChallenge, nil))
}

func (rs *challengeResource) Delete(w http.ResponseWriter, r *http.Request) {
	challenge, ok := r.Context().Value(challengeContextKey).(*db.Challenge)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

	// Only admin or owner can delete the challenge
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

	if !user.IsAdmin && challenge.UserID != user.ID {
		render.Render(w, r, ErrUnauthorized())
		return
	}

	err = rs.store.DeleteChallengeByID(r.Context(), challenge.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewChallengeResponse(challenge, nil))
}
