package service

import (
	"auth/internal/infra"
	"auth/internal/repository"
	"auth/internal/util"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type AuthService interface {
	Use(chi.Router)
	Register(http.ResponseWriter, *http.Request) error
	Login(http.ResponseWriter, *http.Request) error
	Refresh(http.ResponseWriter, *http.Request) error
	RequestEmailVerification(http.ResponseWriter, *http.Request) error
}

type authService struct {
	jwtKey    string
	userRepo  repository.UserRepository
	eventRepo repository.EventRepository
	codeRepo  repository.CodeRepository
	email     infra.EmailClient
}

func NewAuthService(
	jwtKey string,
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
	codeRepo repository.CodeRepository,
	email infra.EmailClient,
) AuthService {
	s := &authService{
		jwtKey:    jwtKey,
		userRepo:  userRepo,
		eventRepo: eventRepo,
		codeRepo:  codeRepo,
		email:     email,
	}
	return s
}

func (s *authService) Use(router chi.Router) {
	router.Post("/register", util.E(s.Register))
	router.Post("/login", util.E(s.Login))
	router.Post("/refresh", util.E(s.Refresh))
	router.Post("/email/verify/request", util.E(s.RequestEmailVerification))
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
	Email      string `json:"email" validate:"required,email"`
	Username   string `json:"username" validate:"required,min=2,max=50"`
	Password   string `json:"password" validate:"required,min=8,max=100"`
	VerifyCode string `json:"verify_code" validate:"required,numeric,len=6"`
}

type respRegister struct {
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *authService) Register(w http.ResponseWriter, r *http.Request) error {
	req, err := util.Body[reqRegister](r)
	if err != nil {
		return err
	}

	if !s.codeRepo.CheckEmailVerifyCode(req.Email, req.VerifyCode) {
		return util.BadRequest("invalid verification code")
	}

	hashedPassword, err := util.GenerateHash(req.Password)
	if err != nil {
		return util.InternalServerError("failed to generate password hash")
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
		return util.InternalServerError("failed to save user")
	}

	err = generateJwtToken(s.jwtKey, w, &userClaims{
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
	if err != nil {
		return util.InternalServerError("failed to generate JWT token")
	}

	s.eventRepo.Save(&repository.Event{
		UserID:    &user.ID,
		Action:    repository.EventRegister,
		Detail:    "{}",
		CreatedAt: time.Now(),
	})

	return util.RespondJson(w, respRegister{
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
	req, err := util.Body[reqLogin](r)
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
		return util.NotFound("user not found")
	}

	v, err := util.ValidateHash(user.Password, req.Password)
	if !v.Valid || err != nil {
		return util.Unauthorized("invalid credentials")
	}

	err = generateJwtToken(s.jwtKey, w, &userClaims{
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
	if err != nil {
		return util.InternalServerError("failed to generate JWT token")
	}

	if v.Obsolete {
		s.updateObsoletePassword(user, req.Password)
	}

	user.LastLogin = time.Now()
	s.userRepo.UpdateLastLogin(user)

	s.eventRepo.Save(&repository.Event{
		UserID:    &user.ID,
		Action:    repository.EventLogin,
		Detail:    "{}",
		CreatedAt: time.Now(),
	})

	return util.RespondJson(w, respLogin{
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
}

func (s *authService) Refresh(w http.ResponseWriter, r *http.Request) error {
	claims, err := authenticate(s.jwtKey, r)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByUsername(claims.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return util.NotFound("user not found")
	}

	err = generateJwtToken(s.jwtKey, w, &userClaims{
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
	if err != nil {
		return err
	}

	user.LastLogin = time.Now()
	s.userRepo.UpdateLastLogin(user)

	return nil
}

type reqEmailCode struct {
	Email string `json:"email" validate:"required,email"`
}

func (s *authService) RequestEmailVerification(w http.ResponseWriter, r *http.Request) error {
	req, err := util.Body[reqEmailCode](r)
	if err != nil {
		return err
	}

	code, err := s.codeRepo.SetEmailVerifyCode(req.Email)
	if err != nil {
		return util.InternalServerError("failed to create verification code")
	}

	err = s.email.SendVerifyEmail(req.Email, code)
	if err != nil {
		return util.InternalServerError("failed to send verification email")
	}

	s.eventRepo.Save(&repository.Event{
		UserID:    nil,
		Action:    repository.EventEmail,
		Detail:    "{}",
		CreatedAt: time.Now(),
	})

	return util.RespondJson(w, "verification email sent")
}
