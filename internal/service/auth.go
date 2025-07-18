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
	userRepo  repository.UserRepository
	eventRepo repository.EventRepository
	codeRepo  repository.CodeRepository
	email     infra.EmailClient
}

func NewAuthService(
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
	codeRepo repository.CodeRepository,
	email infra.EmailClient,
) AuthService {
	s := &authService{
		userRepo:  userRepo,
		eventRepo: eventRepo,
		codeRepo:  codeRepo,
		email:     email,
	}
	return s
}

func (s *authService) Use(router chi.Router) {
	router.Post("/register", util.EH(s.Register))
	router.Post("/login", util.EH(s.Login))
	router.Post("/refresh", util.EH(s.Refresh))
	router.Post("/email/verify/request", util.EH(s.RequestEmailVerification))
}

type reqRegister struct {
	App        string `json:"app" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Username   string `json:"username" validate:"required,min=2,max=16"`
	Password   string `json:"password" validate:"required,min=8,max=100"`
	VerifyCode string `json:"verify_code" validate:"required,numeric,len=6"`
}

func (s *authService) Register(w http.ResponseWriter, r *http.Request) error {
	req, err := util.Body[reqRegister](r)
	if err != nil {
		return err
	}
	if err := util.ValidUsername(req.Username); err != nil {
		return err
	}
	if err := util.ValidPassword(req.Password); err != nil {
		return err
	}
	if !s.codeRepo.CheckEmailVerifyCode(req.Email, req.VerifyCode) {
		return util.BadRequest("无效验证码")
	}

	hashedPassword, err := util.GenerateHash(req.Password)
	if err != nil {
		return util.InternalServerError("密码哈希失败")
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
		return util.InternalServerError("用户保存失败")
	}

	s.eventRepo.Save(&repository.Event{
		UserID:    &user.ID,
		Action:    repository.EventRegister,
		Detail:    "{}",
		CreatedAt: time.Now(),
	})

	return util.RespondAuthTokens(w, util.TokenOptions{
		App:              req.App,
		Username:         user.Username,
		Role:             user.Role,
		CreatedAt:        user.CreatedAt,
		WithRefreshToken: true,
	})
}

type reqLogin struct {
	App      string `json:"app" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *authService) Login(w http.ResponseWriter, r *http.Request) error {
	req, err := util.Body[reqLogin](r)
	if err != nil {
		return err
	}

	var user *repository.User
	if strings.Contains(req.Username, "@") {
		user, err = s.userRepo.FindByEmail(req.Username)
	}
	if err != nil {
		user, err = s.userRepo.FindByUsername(req.Username)
	}
	if err != nil {
		return util.NotFound("用户不存在")
	}

	v, err := util.ValidateHash(user.Password, req.Password)
	if !v.Valid || err != nil {
		return util.Unauthorized("密码错误")
	}
	if v.Obsolete {
		newHashedPassword, err := util.GenerateHash(req.Password)
		if err == nil {
			user.Password = newHashedPassword
			s.userRepo.UpdateHashedPassword(user)
		}
	}

	user.LastLogin = time.Now()
	s.userRepo.UpdateLastLogin(user)

	s.eventRepo.Save(&repository.Event{
		UserID:    &user.ID,
		Action:    repository.EventLogin,
		Detail:    "{}",
		CreatedAt: time.Now(),
	})
	return util.RespondAuthTokens(w, util.TokenOptions{
		App:              req.App,
		Username:         user.Username,
		Role:             user.Role,
		CreatedAt:        user.CreatedAt,
		WithRefreshToken: true,
	})
}

func (s *authService) Refresh(w http.ResponseWriter, r *http.Request) error {
	username, err := util.VerifyRefreshToken(r)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return err
	}
	if user == nil {
		return util.NotFound("用户不存在")
	}

	user.LastLogin = time.Now()
	s.userRepo.UpdateLastLogin(user)

	return util.RespondAuthTokens(w, util.TokenOptions{
		App:              r.URL.Query().Get("app"),
		Username:         user.Username,
		Role:             user.Role,
		CreatedAt:        user.CreatedAt,
		WithRefreshToken: false,
	})
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
		return util.InternalServerError("创建验证码失败")
	}

	err = s.email.SendVerifyEmail(req.Email, code)
	if err != nil {
		return util.InternalServerError("发送验证邮件失败")
	}

	s.eventRepo.Save(&repository.Event{
		UserID:    nil,
		Action:    repository.EventEmail,
		Detail:    "{}",
		CreatedAt: time.Now(),
	})

	return util.RespondJson(w, "验证邮件已发送")
}
