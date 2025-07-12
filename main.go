package main

import (
	"auth/internal/infra"
	"auth/internal/repository"
	"auth/internal/service"
	"auth/internal/util"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return fallback
}

func main() {
	// util
	util.JwtKey = env("JWT_SECRET", "secret")

	// infra
	db := infra.NewSqlDb(
		env("DB_HOST", "localhost"),
		envInt("DB_PORT", 5432),
		env("DB_USER", "auth"),
		env("DB_PASSWORD", ""),
		env("DB_NAME", "auth"),
	)
	rdb := infra.NewRedis(
		env("RDB_HOST", "localhost"),
		envInt("RDB_PORT", 6379),
		env("RDB_USER", "auth"),
		env("RDB_PASSWORD", ""),
	)
	email := infra.NewEmailClient(
		env("MAILGUN_DOMAIN", ""),
		env("MAILGUN_APIKEY", ""),
	)

	// repository
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	codeRepo := repository.NewCodeRepository(rdb)

	// service
	authService := service.NewAuthService(
		userRepo,
		eventRepo,
		codeRepo,
		email,
	)
	adminService := service.NewAdminService(
		userRepo,
		eventRepo,
	)

	// router
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	router.Route("/api/v1", func(router chi.Router) {
		router.Route("/auth", authService.Use)
		router.Route("/admin", adminService.Use)
	})
	http.ListenAndServe(":3000", router)
}
