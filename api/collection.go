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

type collectionResource struct {
	store *db.Store
	auth  *auth.Auth
}

func NewCollectionResource(s *db.Store, a *auth.Auth) *collectionResource {
	return &collectionResource{store: s, auth: a}
}

func (rs *collectionResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(rs.auth.Authenticator())

	r.Get("/", rs.List)
	r.Get("/active", rs.ListActive)
	r.Post("/", rs.Create)

	r.Route("/{collectionID}", func(r chi.Router) {
		r.Use(rs.CollectionCtx)
		r.Get("/", rs.FindOne)
		r.Delete("/", rs.Delete)
		r.Put("/", rs.Update)
	})

	return r
}

type CollectionRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	GoalPoints   int    `json:"goalPoints"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	CollectionID int    `json:"collectionId"`
}

func (c *CollectionRequest) Bind(r *http.Request) error {
	c.Title = strings.TrimSpace(c.Title)
	c.Description = strings.TrimSpace(c.Description)

	if c.Title == "" {
		return errors.New("kokelman nimi on pakollinen")
	}
	if utf8.RuneCountInString(c.Title) > 50 {
		return errors.New("kokelman nimi on liian pitkä (max 50)")
	}

	if c.Description == "" {
		return errors.New("kokelman kuvaus on pakollinen")
	}
	if utf8.RuneCountInString(c.Description) > 2000 {
		return errors.New("kokelman kuvaus on liian pitkä (max 2000)")
	}

	if c.GoalPoints < 1 {
		return errors.New("kokelman tavoitepisteet tulee olla yli 1")
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

type CollectionResponse struct {
	*db.Collection
	Pot int `json:"pot"`
}

func NewCollectionResponse(collection *db.Collection) *CollectionResponse {
	return &CollectionResponse{Collection: collection}
}

func (c *CollectionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	for _, ch := range c.Challenges {
		for _, p := range ch.Performances {
			c.Pot += p.Points
		}
	}
	return nil
}

func NewCollectionListResponse(collections []*db.Collection) []render.Renderer {
	list := []render.Renderer{}
	for _, c := range collections {
		list = append(list, NewCollectionResponse(c))
	}
	return list
}

func (rs *collectionResource) CollectionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		collectionIDStr := chi.URLParam(r, "collectionID")
		if collectionIDStr == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("collection id is missing")))
			return
		}

		collectionID, err := strconv.Atoi(collectionIDStr)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		collection, err := rs.store.FindCollectionByID(r.Context(), collectionID)
		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), collectionContextKey, collection)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *collectionResource) List(w http.ResponseWriter, r *http.Request) {
	collections, err := rs.store.ListCollections(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.RenderList(w, r, NewCollectionListResponse(collections))
}

func (rs *collectionResource) ListActive(w http.ResponseWriter, r *http.Request) {
	collections, err := rs.store.ListActiveCollections(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.RenderList(w, r, NewCollectionListResponse(collections))
}

func (rs *collectionResource) Create(w http.ResponseWriter, r *http.Request) {
	data := &CollectionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userID, ok := getUserID(r)
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no valid user id")))
		return
	}

	newCollection := &db.Collection{
		Title:       data.Title,
		Description: data.Description,
		GoalPoints:  data.GoalPoints,
		StartDate:   data.StartDate,
		EndDate:     data.EndDate,
		UserID:      userID,
	}

	if err := rs.store.AddCollection(r.Context(), newCollection); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewCollectionResponse(newCollection))
}

func (rs *collectionResource) FindOne(w http.ResponseWriter, r *http.Request) {
	collection, ok := r.Context().Value(collectionContextKey).(*db.Collection)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

	render.Render(w, r, NewCollectionResponse(collection))
}

func (rs *collectionResource) Update(w http.ResponseWriter, r *http.Request) {
	currentCollection, ok := r.Context().Value(collectionContextKey).(*db.Collection)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

	data := &CollectionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	currentCollection.Title = data.Title
	currentCollection.Description = data.Description
	currentCollection.GoalPoints = data.GoalPoints
	currentCollection.StartDate = data.StartDate
	currentCollection.EndDate = data.EndDate

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

	if !user.IsAdmin && currentCollection.UserID != user.ID {
		render.Render(w, r, ErrUnauthorized())
		return
	}

	if err := rs.store.UpdateCollection(r.Context(), currentCollection); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewCollectionResponse(currentCollection))
}

func (rs *collectionResource) Delete(w http.ResponseWriter, r *http.Request) {
	collection, ok := r.Context().Value(collectionContextKey).(*db.Collection)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

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

	if !user.IsAdmin && collection.UserID != user.ID {
		render.Render(w, r, ErrUnauthorized())
		return
	}

	err = rs.store.DeleteCollectionByID(r.Context(), collection.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewCollectionResponse(collection))
}
