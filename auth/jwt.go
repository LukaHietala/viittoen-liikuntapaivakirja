package auth

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/db"
)

type Auth struct {
	JWTAuth *jwtauth.JWTAuth
	store   *db.Store
}

func New(secret string, store *db.Store) *Auth {
	return &Auth{
		JWTAuth: jwtauth.New("HS256", []byte(secret), nil, jwt.WithAcceptableSkew(30*time.Second)),
		store:   store,
	}
}

func (a *Auth) MakeSessionToken(userID int) (string, error) {
	claims := map[string]any{
		"user_id": userID,
	}
	jwtauth.SetExpiryIn(claims, 24*28*time.Hour)
	jwtauth.SetIssuedNow(claims)

	_, tokenStr, err := a.JWTAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (a *Auth) MakeResetToken(email string) (string, error) {
	claims := map[string]any{
		"email": email,
	}
	jwtauth.SetExpiryIn(claims, 24*time.Hour)
	jwtauth.SetIssuedNow(claims)

	_, tokenStr, err := a.JWTAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (a *Auth) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(a.JWTAuth)
}

func (a *Auth) Authenticator() func(http.Handler) http.Handler {
	return jwtauth.Authenticator(a.JWTAuth)
}
