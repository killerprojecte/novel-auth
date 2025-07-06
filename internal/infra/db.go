package infra

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"
)

func NewSqlDb(
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

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// err := rdb.Set("key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
}
