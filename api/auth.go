package api

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func InitAuth(secret string) {
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func MakeSessionToken(userId int) string {
	claims := map[string]any{
		"user_id": userId,
	}
	// TODO: refresh tokens maybe?
	jwtauth.SetExpiryIn(claims, 24*31*time.Hour)
	jwtauth.SetIssuedNow(claims)
	_, token, _ := tokenAuth.Encode(claims)
	return token
}

func MakeResetToken(email string) string {
	claims := map[string]any{
		"email": email,
	}
	jwtauth.SetExpiryIn(claims, 15*time.Minute)
	jwtauth.SetIssuedNow(claims)
	_, token, _ := tokenAuth.Encode(claims)
	return token
}