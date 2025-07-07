package service

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type userClaims struct {
	Username  string
	Role      string
	CreatedAt time.Time
}

type userClaimsRaw struct {
	jwt.StandardClaims
	Role      string `json:"role"`
	CreatedAt int64  `json:"crat"`
}

func generateJwtToken(jwtKey string, w http.ResponseWriter, user *userClaims) error {
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(time.Hour * 24 * 30)

	claims := &userClaimsRaw{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: expiredAt.Unix(),
			IssuedAt:  issuedAt.Unix(),
		},
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Unix(),
	}

	token, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(jwtKey)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "auth",
		Value:    token,
		Domain:   ".novelia.cc",
		Path:     "/",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return nil
}

func authenticate(jwtKey string, r *http.Request) (userClaims, error) {
	var zero userClaims

	cookie, err := r.Cookie("auth")
	if err != nil {
		return zero, unauthorized("missing authentication token")
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &userClaimsRaw{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil || !token.Valid {
		return zero, unauthorized("invalid authentication token")
	}

	claims, ok := token.Claims.(*userClaimsRaw)
	if !ok {
		return zero, unauthorized("invalid authentication token")
	}

	return userClaims{
		Username:  claims.Subject,
		Role:      claims.Role,
		CreatedAt: time.Unix(claims.CreatedAt, 0),
	}, nil
}
