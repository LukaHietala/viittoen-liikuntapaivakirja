package api

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/auth"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/db"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/services"
)

type userResource struct {
	store *db.Store
	auth  *auth.Auth
	mail  *services.MailService
}

func NewUserResource(s *db.Store, a *auth.Auth, m *services.MailService) *userResource {
	return &userResource{store: s, auth: a, mail: m}
}

func (rs *userResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(rs.auth.Authenticator())

	r.Get("/", rs.List)
	r.Post("/", rs.Create)

	r.Route("/{userID}", func(r chi.Router) {
		r.Use(rs.UserCtx)
		r.Put("/", rs.Update)
		r.Delete("/", rs.Delete)
	})

	return r
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
	Color    string `json:"color"`
}

func (ur *UserRequest) Bind(r *http.Request) error {
	ur.Name = strings.TrimSpace(ur.Name)
	ur.Email = strings.TrimSpace(ur.Email)
	ur.Password = strings.TrimSpace(ur.Password)

	if ur.Name == "" {
		return errors.New("käyttäjän nimi on pakollinen")
	}
	if utf8.RuneCountInString(ur.Name) > 50 {
		return errors.New("käyttäjän nimi on liian pitkä")
	}
	if ur.Color == "" {
		return errors.New("käyttäjän väri on pakollinen")
	}
	if ur.Email == "" {
		return errors.New("käyttäjän sähköposti on pakollinen")
	}
	_, err := mail.ParseAddress(ur.Email)
	if err != nil {
		return errors.New("käyttäjän sähköposti ei ole oikeassa muodossa")
	}

	return nil
}

type UserResponse struct {
	*db.User
}

func NewUserResponse(user *db.User) *UserResponse {
	return &UserResponse{User: user}
}

func (ur *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	ur.PasswordHash = ""
	return nil
}

func NewUserListResponse(users []*db.User) []render.Renderer {
	list := []render.Renderer{}
	for _, user := range users {
		list = append(list, NewUserResponse(user))
	}
	return list
}

func (rs *userResource) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "userID")
		if userIDStr == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("user id is missing")))
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		user, err := rs.store.FindUserByID(r.Context(), userID)
		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *userResource) List(w http.ResponseWriter, r *http.Request) {
	users, err := rs.store.ListUsers(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.RenderList(w, r, NewUserListResponse(users))
}

func (rs *userResource) Create(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newUser := &db.User{
		Name:    data.Name,
		Email:   data.Email,
		Color:   data.Color,
		IsAdmin: data.IsAdmin,
	}

	length := 5
	b := make([]byte, length+2)
	rand.Read(b)
	password := fmt.Sprintf("%x", b)[2 : length+2]
	newUser.SetPasswordHash(password)

	if err := rs.store.AddUser(r.Context(), newUser); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	body := fmt.Sprintf(
		"Sinulle on luotu uusi käyttäjä Viittoenliikkuu sovellukseen\n\nNimi: %s\nSalasana: %s\n", data.Name, password,
	)
	if err := rs.mail.Send("Tunnukset viittoenliikkuu.fi", body, data.Email); err != nil {
		// TODO: VALITA
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserResponse(newUser))
}

func (rs *userResource) Update(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value(userContextKey).(*db.User)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	currentUser.Name = data.Name
	currentUser.Email = data.Email
	currentUser.IsAdmin = data.IsAdmin
	currentUser.Color = data.Color

	if data.Password != "" {
		currentUser.SetPasswordHash(data.Password)
	}

	if err := rs.store.UpdateUser(r.Context(), currentUser); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewUserResponse(currentUser))
}

func (rs *userResource) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(*db.User)
	if !ok {
		render.Render(w, r, ErrInternal(errors.New("context error")))
		return
	}

	err := rs.store.DeleteUserByID(r.Context(), user.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewUserResponse(user))
}
