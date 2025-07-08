package main

import (
	"auth/internal/infra"
	"auth/internal/repository"
	"auth/internal/service"
	"log"
	"net/http"
	"os"
	"strconv"
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
	mux := http.NewServeMux()

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

	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	codeRepo := repository.NewCodeRepository(rdb)

	jwt_secret := env("JWT_SECRET", "secret")
	authService := service.NewAuthService(
		jwt_secret,
		userRepo,
		eventRepo,
		codeRepo,
		email,
	)
	adminService := service.NewAdminService(
		jwt_secret,
		userRepo,
		eventRepo,
	)

	const root = "/api/v1"
	service.UseAuthService(mux, authService, root+"/auth")
	service.UseAdminService(mux, adminService, root+"/admin")

	log.Print("Listening... http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
