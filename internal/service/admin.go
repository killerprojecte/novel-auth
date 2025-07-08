package service

import (
	"auth/internal/repository"
	"net/http"
)

func UseAdminService(mux *http.ServeMux, s AdminService, path string) {
	mux.HandleFunc("DELETE "+path+"/user", s.DeleteUser)
	mux.HandleFunc("POST "+path+"/user/role", s.UpdateUserRole)
}

type AdminService interface {
	DeleteUser(http.ResponseWriter, *http.Request)
	UpdateUserRole(http.ResponseWriter, *http.Request)
}

type adminService struct {
	jwtKey    string
	userRepo  repository.UserRepository
	eventRepo repository.EventRepository
}

func NewAdminService(
	jwtKey string,
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
) AdminService {
	s := &adminService{
		jwtKey:    jwtKey,
		userRepo:  userRepo,
		eventRepo: eventRepo,
	}
	return s
}

func (s *adminService) DeleteUser(w http.ResponseWriter, r *http.Request) {
}

func (s *adminService) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
}
