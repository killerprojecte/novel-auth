package service

import (
	"auth/internal/infra"
	"auth/internal/repository"
	"auth/internal/util"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	EventLogin         string = "login"
	EventRegister      string = "register"
	EventLogout        string = "logout"
	EventOtp           string = "otp"
	EventResetPassword string = "reset_password"
)

type AuthService interface {
	Use(chi.Router)
	Register(http.ResponseWriter, *http.Request) error
	Login(http.ResponseWriter, *http.Request) error
	Refresh(http.ResponseWriter, *http.Request) error
	Logout(http.ResponseWriter, *http.Request) error
	RequestOtp(http.ResponseWriter, *http.Request) error
	ResetPassword(http.ResponseWriter, *http.Request) error
}

type authService struct {
	userRepo  repository.UserRepository
	eventRepo repository.EventRepository
	otpRepo   repository.OtpRepository
	email     infra.EmailClient
}

func NewAuthService(
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
	otpRepo repository.OtpRepository,
	email infra.EmailClient,
) AuthService {
	s := &authService{
		userRepo:  userRepo,
		eventRepo: eventRepo,
		otpRepo:   otpRepo,
		email:     email,
	}
	return s
}

func (s *authService) Use(router chi.Router) {
	router.Group(func(router chi.Router) {
		router.Use(util.RateLimiter(10))
		router.Post("/register", util.EH(s.Register))
	})
	router.Group(func(router chi.Router) {
		router.Use(util.RateLimiter(20))
		router.Post("/login", util.EH(s.Login))
	})
	router.Group(func(router chi.Router) {
		router.Use(util.RateLimiter(5))
		router.Post("/otp/request", util.EH(s.RequestOtp))
	})
	router.Group(func(router chi.Router) {
		router.Use(util.RateLimiter(100))
		router.Post("/refresh", util.EH(s.Refresh))
	})
	router.Group(func(router chi.Router) {
		router.Use(util.RateLimiter(5))
		router.Post("/password/reset", util.EH(s.ResetPassword))
	})
}

func (s *authService) Register(w http.ResponseWriter, r *http.Request) error {
	req, err := util.Body[struct {
		App      string `json:"app" validate:"required"`
		Username string `json:"username" validate:"required,min=2,max=16"`
		Password string `json:"password" validate:"required,min=8,max=100"`
		Email    string `json:"email" validate:"required,email"`
		Otp      string `json:"otp" validate:"required,numeric,len=6"`
	}](r)
	if err != nil {
		return err
	}
	if err := util.ValidUsername(req.Username); err != nil {
		return err
	}
	if err := util.ValidPassword(req.Password); err != nil {
		return err
	}
	if !s.otpRepo.CheckOtp(repository.OtpVerify, req.Email, req.Otp) {
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

	s.eventRepo.Save(
		EventRegister,
		&struct {
			App        string `json:"app"`
			ActorUser  string `json:"actor_user"`
			TargetUser string `json:"target_user"`
		}{
			App:        req.App,
			ActorUser:  user.Username,
			TargetUser: user.Username,
		},
	)

	return util.RespondAuthTokens(w, util.TokenOptions{
		App:              req.App,
		Username:         user.Username,
		Role:             user.Role,
		CreatedAt:        user.CreatedAt,
		WithRefreshToken: true,
	})
}

func (s *authService) Login(w http.ResponseWriter, r *http.Request) error {
	req, err := util.Body[struct {
		App      string `json:"app" validate:"required"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}](r)
	if err != nil {
		return err
	}

	var user *repository.User
	if strings.Contains(req.Username, "@") {
		user, err = s.userRepo.FindByEmail(req.Username)
		if err != nil {
			return util.InternalServerError("查询用户失败")
		}
	}
	if user == nil {
		user, err = s.userRepo.FindByUsername(req.Username)
		if err != nil {
			return util.InternalServerError("查询用户失败")
		}
	}
	if user == nil {
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

	s.eventRepo.Save(
		EventLogin,
		&struct {
			App        string `json:"app"`
			ActorUser  string `json:"actor_user"`
			TargetUser string `json:"target_user"`
		}{
			App:        req.App,
			ActorUser:  user.Username,
			TargetUser: user.Username,
		},
	)
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

func (s *authService) Logout(w http.ResponseWriter, r *http.Request) error {
	username, err := util.VerifyRefreshToken(r)
	if err != nil {
		return err
	}

	s.eventRepo.Save(
		EventLogout,
		&struct {
			ActorUser  string `json:"actor_user"`
			TargetUser string `json:"target_user"`
		}{
			ActorUser:  username,
			TargetUser: username,
		},
	)

	return util.RespondLogout(w)
}

func (s *authService) sendOtpEmail(otpType string, email string, otp string) error {
	switch otpType {
	case repository.OtpVerify:
		return s.email.SendEmail(
			email,
			fmt.Sprintf(
				"%s 轻小说机翻机器人 注册激活码",
				otp,
			),
			fmt.Sprintf(
				"您的注册激活码为 %s\n"+
					"激活码将会在15分钟后失效,请尽快完成注册\n"+
					"这是系统邮件，请勿回复",
				otp,
			),
		)
	case repository.OtpResetPassword:
		return s.email.SendEmail(
			email,
			fmt.Sprintf(
				"%s 轻小说机翻机器人 重置密码验证码",
				otp,
			),
			fmt.Sprintf(
				"您的重置密码验证码为 %s\n"+
					"验证码将会在15分钟后失效,请尽快完成操作\n"+
					"这是系统邮件，请勿回复",
				otp,
			),
		)
	default:
		return fmt.Errorf("未知的Otp类型: %s", otpType)
	}
}

func (s *authService) RequestOtp(w http.ResponseWriter, r *http.Request) error {
	req, err := util.Body[struct {
		Email string `json:"email" validate:"required,email"`
		Type  string `json:"type" validate:"required,oneof=verify reset_password"`
	}](r)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return util.InternalServerError("邮件检查失败")
	}

	// 根据不同类型进行不同的验证
	switch req.Type {
	case repository.OtpVerify:
		if user != nil {
			return util.Conflict("邮箱已经被使用")
		}
	case repository.OtpResetPassword:
		if user == nil {
			return util.NotFound("用户不存在")
		}
	default:
		return util.BadRequest("无效的请求类型")
	}

	otp, err := s.otpRepo.SetOtp(req.Type, req.Email)
	if err != nil {
		return util.InternalServerError("创建验证码失败")
	}

	err = s.sendOtpEmail(req.Type, req.Email, otp)
	if err != nil {
		return util.InternalServerError("发送验证邮件失败")
	}

	s.eventRepo.Save(
		EventOtp,
		&struct {
			Email string `json:"email"`
			Type  string `json:"type"`
		}{
			Email: req.Email,
			Type:  req.Type,
		},
	)

	return util.RespondText(w, "验证邮件已发送")
}

func (s *authService) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	req, err := util.Body[struct {
		Email    string `json:"email" validate:"required,email"`
		Otp      string `json:"otp" validate:"required,len=32"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}](r)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return util.InternalServerError("查询用户失败")
	}
	if user == nil {
		return util.NotFound("用户不存在")
	}

	if !s.otpRepo.CheckOtp(repository.OtpResetPassword, req.Email, req.Otp) {
		return util.Unauthorized("无效的验证码")
	}

	newHashedPassword, err := util.GenerateHash(req.Password)
	if err != nil {
		return util.InternalServerError("密码哈希失败")
	}
	user.Password = newHashedPassword
	err = s.userRepo.UpdateHashedPassword(user)
	if err != nil {
		return util.InternalServerError("密码重置失败")
	}

	s.eventRepo.Save(
		EventResetPassword,
		&struct {
			ActorUser  string `json:"actor_user"`
			TargetUser string `json:"target_user"`
		}{
			ActorUser:  user.Username,
			TargetUser: user.Username,
		},
	)

	return util.RespondText(w, "密码重置成功")
}
