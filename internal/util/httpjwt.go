package util

import (
	"auth/internal/repository"
	"net/http"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type VerifyLevel string

const (
	LevelAdmin  VerifyLevel = "admin"
	LevelMember VerifyLevel = "member"
)

type VerifiedUser struct {
	Username  string
	Role      string
	CreatedAt time.Time
	ExpiredAt time.Time
}

var JwtKey string

func GetVerifiedUser(r *http.Request, level VerifyLevel) (VerifiedUser, error) {
	claims, err := getClaimFromCookie(r)
	if err != nil {
		return VerifiedUser{}, err
	}

	switch level {
	case LevelAdmin:
		if claims.Role != repository.RoleAdmin {
			return VerifiedUser{}, Unauthorized("insufficient privileges")
		}
		// 危险接口限制 token 有效时间
		if time.Since(claims.IssuedAt.Time).Minutes() > 15 {
			return VerifiedUser{}, Unauthorized("user is banned")
		}
	case LevelMember:
		validRoles := []string{
			repository.RoleAdmin,
			repository.RoleMember,
			repository.RoleGuest,
		}
		if !slices.Contains(validRoles, claims.Role) {
			return VerifiedUser{}, Unauthorized("insufficient privileges")
		}
	default:
		return VerifiedUser{}, Unauthorized("invalid verification level")
	}

	if level == LevelMember {
		if claims.Role != string(LevelMember) && claims.Role != string(LevelAdmin) {
			return VerifiedUser{}, Unauthorized("member privileges required")
		}
	}
	if level == LevelAdmin && claims.Role != string(LevelAdmin) {
		return VerifiedUser{}, Unauthorized("admin privileges required")
	}

	return VerifiedUser{
		Username:  claims.Subject,
		Role:      claims.Role,
		CreatedAt: claims.CreatedAt.Time,
	}, nil
}

func IssueToken(w http.ResponseWriter, user *VerifiedUser) error {
	if err := setClaimToCookie(w, user); err != nil {
		return err
	}
	return nil
}

func RefreshToken(w http.ResponseWriter, user *VerifiedUser) error {
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

func setClaimToCookie(w http.ResponseWriter, user *VerifiedUser) error {
	issuedAt := time.Now()
	expiredAt := user.ExpiredAt
	if expiredAt.IsZero() {
		expiredAt = issuedAt.Add(time.Hour * 24 * 30)
	}

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
		MaxAge:   int(expiredAt.Sub(issuedAt).Seconds() - 60),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return nil
}
