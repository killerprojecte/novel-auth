package util

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthenticateLevel string

const (
	LevelUser  AuthenticateLevel = "user"
	LevelAdmin AuthenticateLevel = "admin"
)

type VerifiedUser struct {
	Username  string
	Role      string
	CreatedAt time.Time
}

type LoginUser = VerifiedUser

var JwtKey string

func GetVerifiedUser(r *http.Request) (VerifiedUser, error) {
	claims, err := getClaimFromCookie(r)
	if err != nil {
		return VerifiedUser{}, err
	}
	return VerifiedUser{
		Username:  claims.Subject,
		Role:      claims.Role,
		CreatedAt: claims.CreatedAt.Time,
	}, nil
}

func IssueToken(w http.ResponseWriter, user *LoginUser) error {
	if err := setClaimToCookie(w, user); err != nil {
		return err
	}
	return nil
}

func RefreshToken(w http.ResponseWriter, user *LoginUser, verifiedUser *VerifiedUser) error {
	if err := setClaimToCookie(w, user); err != nil {
		return err
	}
	return nil
}

type userClaim struct {
	jwt.RegisteredClaims
	Role      string           `json:"role"`
	CreatedAt *jwt.NumericDate `json:"crat"`
}

func getClaimFromCookie(r *http.Request) (userClaim, error) {
	zero := userClaim{}

	cookie, err := r.Cookie("auth")
	if err != nil {
		return zero, Unauthorized("missing authentication token")
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &userClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		},
	)
	if err != nil || !token.Valid {
		return zero, Unauthorized("missing authentication token")
	}

	claims, ok := token.Claims.(*userClaim)
	if !ok {
		return zero, Unauthorized("invalid authentication token")
	}
	return *claims, nil
}

func setClaimToCookie(w http.ResponseWriter, user *LoginUser) error {
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(time.Hour * 24 * 30)

	claims := &userClaim{
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
		SignedString(JwtKey)
	if err != nil {
		return InternalServerError("failed to generate JWT token")
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
