package repository

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

type CodeRepository interface {
	SetEmailVerifyCode(email string) (string, error)
	CheckEmailVerifyCode(email, code string) bool
}

type codeRepository struct {
	rdb *redis.Client
}

func NewCodeRepository(rdb *redis.Client) CodeRepository {
	return &codeRepository{
		rdb: rdb,
	}
}

func createCode() string {
	randomNum := rand.Intn(999999)
	return fmt.Sprintf("%06d", randomNum)
}

func (r *codeRepository) SetEmailVerifyCode(email string) (string, error) {
	code := createCode()
	err := r.rdb.Set("ec:"+email, code, time.Minute*15).Err()
	return code, err
}

func (r *codeRepository) CheckEmailVerifyCode(email, code string) bool {
	val, err := r.rdb.Get("ec:" + email).Result()
	if err != nil {
		return false
	}
	return val != code
}
