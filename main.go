package main

import (
	"auth/internal/infra"
	"auth/internal/repository"
	"auth/internal/service"
	"log"
	"net/http"
	"os"
)

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	mux := http.NewServeMux()

	db := infra.NewSqlDb(
		env("DB_HOST", "localhost"),
		env("DB_USER", "auth"),
		env("DB_PASSWORD", ""),
		env("DB_NAME", "auth"),
	)
	rdb := infra.NewRedis()

	userRepo := repository.NewUserRepository(db)
	codeRepo := repository.NewCodeRepository(rdb)
	authService := service.NewAuthService("jwt-secret-key", userRepo, codeRepo)
	adminService := service.NewAdminService(userRepo)

	const root = "/api/v1"
	service.UseAuthService(mux, authService, root+"/auth")
	service.UseAdminService(mux, adminService, root+"/admin")

	log.Print("Listening... http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
