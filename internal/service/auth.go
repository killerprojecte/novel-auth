package service

import (
	"auth/internal/repository"
	"auth/internal/util"
	"net/http"
	"strings"
	"time"
)

func UseAuthService(mux *http.ServeMux, s AuthService, path string) {
	mux.HandleFunc("POST "+path+"/register", ToHandler(s.Register))
	mux.HandleFunc("POST "+path+"/login", ToHandler(s.Login))
	mux.HandleFunc("POST "+path+"/refresh", ToHandler(s.Refresh))
	mux.HandleFunc("POST "+path+"/email/verify/request", ToHandler(s.RequestEmailVerification))
	mux.HandleFunc("POST "+path+"/password/reset", ToHandler(s.ResetPassword))
	mux.HandleFunc("POST "+path+"/password/reset/request", ToHandler(s.RequestPasswordReset))
}

type AuthService interface {
	Register(http.ResponseWriter, *http.Request) error
	Login(http.ResponseWriter, *http.Request) error
	Refresh(http.ResponseWriter, *http.Request) error
	RequestEmailVerification(http.ResponseWriter, *http.Request) error
	ResetPassword(http.ResponseWriter, *http.Request) error
	RequestPasswordReset(http.ResponseWriter, *http.Request) error
}

type authService struct {
	jwtKey   string
	userRepo repository.UserRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
) AuthService {
	s := &authService{
		jwtKey:   "wtf",
		userRepo: userRepo,
	}
	return s
}
func (s *authService) generateAndUseJwtToken(w http.ResponseWriter, user *repository.User) error {
	expired := time.Now().Add(time.Hour * 24 * 30)
	token, err := util.GenerateJwt(s.jwtKey, user, expired)
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

func (s *authService) updateObsoletePassword(user *repository.User, newPassword string) error {
	hashedPassword, err := util.GenerateHash(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.UpdateHashedPassword(user)
}

type reqRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type respRegister struct {
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Role      string    `json:"role" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

func (s *authService) Register(w http.ResponseWriter, r *http.Request) error {
	req, err := body[reqRegister](r)
	if err != nil {
		return err
	}

	hashedPassword, err := util.GenerateHash(req.Password)
	if err != nil {
		return InternalServerError("failed to generate password hash")
	}

	user := &repository.User{
		Username:  req.Username,
		Email:     req.Email,
		Role:      repository.RoleMember,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		LastLogin: time.Now(),
		Attr:      "{}",
	}
	err = s.userRepo.Save(user)
	if err != nil {
		return InternalServerError("failed to save user")
	}

	err = s.generateAndUseJwtToken(w, user)
	if err != nil {
		return InternalServerError("failed to generate JWT token")
	}

	return Respond(w, http.StatusCreated, respRegister{
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
}

type reqLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type respLogin struct {
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Role      string    `json:"role" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

func (s *authService) Login(w http.ResponseWriter, r *http.Request) error {
	req, err := body[reqLogin](r)
	if err != nil {
		return err
	}

	var user *repository.User
	if strings.Contains(req.Username, "@") {
		user, err = s.userRepo.FindByEmail(req.Username)
	} else {
		user, err = s.userRepo.FindByUsername(req.Username)
	}
	if err != nil {
		return NotFound("user not found")
	}

	v, err := util.ValidateHash(user.Password, req.Password)
	if !v.Valid || err != nil {
		return Unauthorized("invalid credentials")
	}

	err = s.generateAndUseJwtToken(w, user)
	if err != nil {
		return InternalServerError("failed to generate JWT token")
	}

	if v.Obsolete {
		err = s.updateObsoletePassword(user, req.Password)
		if err != nil {
		}
	}

	user.LastLogin = time.Now()
	err = s.userRepo.UpdateLastLogin(user)
	if err != nil {
	}

	return Respond(w, http.StatusOK, respLogin{
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
}

func (s *authService) Refresh(w http.ResponseWriter, r *http.Request) error {
	// Implement refresh logic here

	return nil
}

func (s *authService) RequestEmailVerification(w http.ResponseWriter, r *http.Request) error {
	// Implement email verification request logic here

	return nil
}

func (s *authService) RequestPasswordReset(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type reqResetPassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (s *authService) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	_, err := body[reqResetPassword](r)
	if err != nil {
		return err
	}

	return nil
}
