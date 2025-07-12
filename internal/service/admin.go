package service

import (
	"auth/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AdminService interface {
	Use(chi.Router)
	DeleteUser(http.ResponseWriter, *http.Request)
	UpdateUserRole(http.ResponseWriter, *http.Request)
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
	router.Delete("/user", s.DeleteUser)
	router.Post("/user/role", s.UpdateUserRole)
}

func (s *adminService) DeleteUser(w http.ResponseWriter, r *http.Request) {
}

func (s *adminService) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
}
