package infra

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func NewSqlDb(host string, port int, user, password, dbname string) *sql.DB {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		),
	)
	if err != nil {
		panic(err)
	}
	return db
}

func NewRedis(host string, port int, user, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Username: user,
		Password: password,
		DB:       0,
	})
}
