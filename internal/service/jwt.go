package service

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userClaims struct {
	Username  string
	Role      string
	CreatedAt time.Time
}

type userClaimsRaw struct {
	jwt.RegisteredClaims
	Role      string           `json:"role"`
	CreatedAt *jwt.NumericDate `json:"crat"`
}

func generateJwtToken(jwtKey string, w http.ResponseWriter, user *userClaims) error {
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(time.Hour * 24 * 30)

	claims := &userClaimsRaw{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Username,
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
		Role:      user.Role,
		CreatedAt: jwt.NewNumericDate(user.CreatedAt),
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
		CreatedAt: claims.CreatedAt.Time,
	}, nil
}
