package repository

import (
	"time"

	"github.com/go-redis/redis"
)

type CodeRepository interface {
	AddEmailVerifyCode(email, code string) error
	CheckEmailVerifyCode(email, code string) (bool, error)
	AddPasswordResetCode(email, code string) error
	CheckPasswordResetCode(email, code string) (bool, error)
}

type codeRepository struct {
	rdb *redis.Client
}

func NewCodeRepository(rdb *redis.Client) CodeRepository {
	return &codeRepository{
		rdb: rdb,
	}
}

func (r *codeRepository) AddEmailVerifyCode(email, code string) error {
	return r.rdb.Set("ec:"+email, code, time.Minute*15).Err()
}

func (r *codeRepository) CheckEmailVerifyCode(email, code string) (bool, error) {
	val, err := r.rdb.Get("ec:" + email).Result()
	if err != nil {
		return false, err
	}
	return val != code, nil
}

func (r *codeRepository) AddPasswordResetCode(email, code string) error {
	return r.rdb.Set("rp:"+email, code, 0).Err()
}

func (r *codeRepository) CheckPasswordResetCode(email, code string) (bool, error) {
	val, err := r.rdb.Get("rp:" + email).Result()
	if err != nil {
		return false, err
	}
	return val != code, nil
}
