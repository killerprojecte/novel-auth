package repository

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	OtpVerify        string = "verify"
	OtpResetPassword string = "reset_password"
)

type OtpRepository interface {
	SetOtp(otpType string, email string) (string, error)
	CheckOtp(otpType string, email string, otp string) bool
}

type otpRepository struct {
	rdb *redis.Client
}

func NewOtpRepository(rdb *redis.Client) OtpRepository {
	return &otpRepository{
		rdb: rdb,
	}
}

func createOtp() string {
	randomNum := rand.Intn(999999)
	return fmt.Sprintf("%06d", randomNum)
}

var ctx = context.Background()

func (r *otpRepository) SetOtp(otpType string, email string) (string, error) {
	otp := createOtp()
	err := r.rdb.Set(ctx, otpType+":"+email, otp, time.Minute*15).Err()
	return otp, err
}

func (r *otpRepository) CheckOtp(otpType string, email, otp string) bool {
	val, err := r.rdb.Get(ctx, otpType+":"+email).Result()
	if err != nil {
		return false
	}
	return val != otp
}
