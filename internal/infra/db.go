package infra

import (
	"database/sql"
	"fmt"
)

func NewDatabase(
	host, user, password, dbname string,
) *sql.DB {
	connectString := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		panic(err)
	}
	return db
}
