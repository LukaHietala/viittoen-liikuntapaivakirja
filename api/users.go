package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/LukaHietala/viittoen-liikuntapaivakirja/db"
	"github.com/LukaHietala/viittoen-liikuntapaivakirja/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type usersResource struct{}

func (rs usersResource) Routes(ms *services.MailService) chi.Router {
	r := chi.NewRouter()
	r.Use(AdminOnly)

	r.Get("/", rs.List)
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		rs.Create(w, r, ms)
	})

	r.Route("/{id}", func(r chi.Router) {
		r.Use(rs.UserCtx)
		r.Delete("/", rs.Delete)
		r.Patch("/", rs.Update)
	})

	return r
}

func (rs usersResource) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *db.User
		var err error

		userIDStr := chi.URLParam(r, "id")
		if userIDStr == "" {
			render.Render(w, r, ErrNotFound())
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		user, err = store.GetUserByID(userID)
		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs usersResource) List(w http.ResponseWriter, r *http.Request) {
	users, err := store.ListUsers()
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, users)
}

func (rs usersResource) Create(w http.ResponseWriter, r *http.Request, ms *services.MailService) {
	var user db.User
	if err := render.Decode(r, &user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := user.Validate(); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if err := store.AddUser(user, ms); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, 201)
	render.JSON(w, r, user)
}

func (rs usersResource) Update(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*db.User)

	var req db.User
	if err := render.Decode(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	user = &req
	if err := store.UpdateUser(r.Context(), user.ID, req); err != nil {
		render.Render(w, r, ErrInternal(err))
	}

	render.Status(r, 200)
	render.JSON(w, r, user)
}

func (rs usersResource) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*db.User)

	err := store.DeleteUser(r.Context(), user.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, 200)
	render.JSON(w, r, user)
}