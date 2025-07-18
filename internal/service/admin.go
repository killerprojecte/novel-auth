package service

import (
	"auth/internal/repository"
	"auth/internal/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type AdminService interface {
	Use(chi.Router)
	GetUser(http.ResponseWriter, *http.Request) error
	RestrictUser(http.ResponseWriter, *http.Request) error
	BanUser(http.ResponseWriter, *http.Request) error
}

type adminService struct {
	userRepo  repository.UserRepository
	eventRepo repository.EventRepository
}

func NewAdminService(
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
) AdminService {
	s := &adminService{
		userRepo:  userRepo,
		eventRepo: eventRepo,
	}
	return s
}

func (s *adminService) Use(router chi.Router) {
	router.Get("/user", util.EH(s.GetUser))
	router.Post("/user/restrict", util.EH(s.RestrictUser))
	router.Post("/user/ban", util.EH(s.BanUser))
}

func (s *adminService) GetUser(w http.ResponseWriter, r *http.Request) error {
	_, err := util.VerifyAccessToken(r, true)
	if err != nil {
		return err
	}
	users, err := s.userRepo.List(repository.UserFilter{}, 0, 10)
	if err != nil {
		return err
	}
	return util.RespondJson(w, users)
}

func (s *adminService) RestrictUser(w http.ResponseWriter, r *http.Request) error {
	adminUsername, err := util.VerifyAccessToken(r, true)
	if err != nil {
		return err
	}

	req, err := util.Body[struct {
		Username string `json:"username" validate:"required"`
		Reason   string `json:"reason" validate:"required"`
	}](r)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return util.NotFound("用户不存在")
	}
	if user.Role != repository.RoleMember {
		return util.Unauthorized("没有权限对非普通用户进行操作")
	}

	user.Role = repository.RoleRestricted
	err = s.userRepo.UpdateRole(user)
	if err != nil {
		return util.InternalServerError("更新用户角色失败")
	}

	detail, _ := json.Marshal(&struct {
		ActorUser  string `json:"actor_user"`
		TargetUser string `json:"target_user"`
		Reason     string `json:"reason"`
	}{
		ActorUser:  adminUsername,
		TargetUser: user.Username,
		Reason:     req.Reason,
	})
	s.eventRepo.Save(&repository.Event{
		UserID:    &user.ID,
		Action:    repository.EventRestrictUser,
		Detail:    string(detail),
		CreatedAt: time.Now(),
	})

	return nil
}

func (s *adminService) BanUser(w http.ResponseWriter, r *http.Request) error {
	adminUsername, err := util.VerifyAccessToken(r, true)
	if err != nil {
		return err
	}

	req, err := util.Body[struct {
		Username string `json:"username" validate:"required"`
		Reason   string `json:"reason" validate:"required"`
	}](r)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return util.NotFound("用户不存在")
	}
	if user.Role != repository.RoleMember {
		return util.Unauthorized("没有权限对非普通用户进行操作")
	}

	user.Role = repository.RoleBanned
	err = s.userRepo.UpdateRole(user)
	if err != nil {
		return util.InternalServerError("更新用户角色失败")
	}

	detail, _ := json.Marshal(&struct {
		ActorUser  string `json:"actor_user"`
		TargetUser string `json:"target_user"`
		Reason     string `json:"reason"`
	}{
		ActorUser:  adminUsername,
		TargetUser: user.Username,
		Reason:     req.Reason,
	})
	s.eventRepo.Save(&repository.Event{
		UserID:    &user.ID,
		Action:    repository.EventBanUser,
		Detail:    string(detail),
		CreatedAt: time.Now(),
	})

	return nil
}
