package api

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/auth"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/db"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/services"
)

type authResource struct {
	store *db.Store
	auth  *auth.Auth
	mail  *services.MailService
}

func NewAuthResource(s *db.Store, a *auth.Auth, m *services.MailService) *authResource {
	return &authResource{store: s, auth: a, mail: m}
}

func (rs *authResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(rs.auth.RejectAuthenticated)
		r.Post("/login", rs.Login)
	})

	r.Post("/logout", rs.Logout)
	r.Get("/self", rs.Self)
	r.Post("/forgot", rs.ForgotPassword)
	r.Post("/reset", rs.ResetPassword)

	return r
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (l *LoginRequest) Bind(r *http.Request) error {
	l.Name = strings.TrimSpace(l.Name)
	l.Password = strings.TrimSpace(l.Password)

	if l.Name == "" {
		return errors.New("nimi on pakollinen")
	}
	if l.Password == "" {
		return errors.New("salasana on pakollinen")
	}

	return nil
}

func (rs *authResource) Login(w http.ResponseWriter, r *http.Request) {
	data := &LoginRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := rs.store.FindUserByName(r.Context(), data.Name)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("käyttäjää ei ole olemassa")))
		return
	}

	if ok := user.Verify(data.Password); !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("nimi tai salasana on väärä")))
		return
	}

	token, err := rs.auth.MakeSessionToken(user.ID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * 28 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		// Secure: true,
		Name:  "jwt",
		Value: token,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (rs *authResource) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		// Secure: true,
		Name:  "jwt",
		Value: "",
	})

	w.WriteHeader(http.StatusNoContent)
}

type SelfResponse struct {
	User          *db.User `json:"user,omitempty"`
	Authenticated bool     `json:"authenticated"`
	// ?
	ToOnboarding bool `json:"toOnboarding"`
}

func NewSelfResponse(user *db.User, authenticated bool, toOnboarding bool) *SelfResponse {
	return &SelfResponse{User: user, Authenticated: authenticated, ToOnboarding: toOnboarding}
}

func (l *SelfResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if l.User != nil {
		l.User.PasswordHash = ""
	}
	return nil
}

func (rs *authResource) Self(w http.ResponseWriter, r *http.Request) {
	userCount, err := rs.store.UserCount()
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	toOnboarding := false
	if userCount == 0 {
		toOnboarding = true
	}

	userID, ok := getUserID(r)
	if !ok {
		render.Render(w, r, NewSelfResponse(nil, false, toOnboarding))
		return
	}

	user, err := rs.store.FindUserByID(r.Context(), userID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Render(w, r, NewSelfResponse(user, true, false))
}

type ForgotRequest struct {
	Email string `json:"email"`
}

func (f *ForgotRequest) Bind(r *http.Request) error {
	f.Email = strings.TrimSpace(f.Email)

	if f.Email == "" {
		return errors.New("sähköposti on pakollinen")
	}
	_, err := mail.ParseAddress(f.Email)
	if err != nil {
		return errors.New("sähköposti ei ole oikeassa muodossa")
	}

	return nil
}

func (rs *authResource) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	data := &ForgotRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	token, err := rs.auth.MakeResetToken(data.Email)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	user, err := rs.store.FindUserByEmail(r.Context(), data.Email)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("käyttäjää ei löytynyt")))
		return
	}

	resetLink := "https://viittoenliikkuu.fi/reset-password?token=" + token

	err = rs.mail.Send("Nollaa salasana", resetLink, user.Email)

	// TODO: Send email
	w.WriteHeader(http.StatusNoContent)
}

type ResetRequest struct {
	NewPassword string `json:"newPassword"`
	Token       string `json:"token"`
}

func (re *ResetRequest) Bind(r *http.Request) error {
	re.NewPassword = strings.TrimSpace(re.NewPassword)

	if re.NewPassword == "" {
		return errors.New("uusi salasana on pakollinen")
	}
	if re.Token == "" {
		return errors.New("linkki on virheellinen")
	}

	// Token validation is done locally

	return nil
}

func (rs *authResource) ResetPassword(w http.ResponseWriter, r *http.Request) {
	data := &ResetRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_, err := rs.auth.JWTAuth.Decode(data.Token)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("linkki on virheellinen")))
		return
	}

	token, err := jwtauth.VerifyToken(rs.auth.JWTAuth, data.Token)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("linkki on vanhentunut")))
		return
	}

	var email string

	err = token.Get(`email`, &email)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("linkki on virheellinen")))
		return
	}

	user, err := rs.store.FindUserByEmail(r.Context(), email)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("käyttäjää ei löytynyt")))
		return
	}

	if err := user.SetPasswordHash(data.NewPassword); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	if err := rs.store.UpdateUser(r.Context(), user); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewUserResponse(user))
}
