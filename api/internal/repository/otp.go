package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
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

func createOtp(otpType string) (string, error) {
	switch otpType {
	case OtpVerify:
		n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%06d", n.Int64()), nil
	case OtpResetPassword:
		return rand.Text(), nil
	default:
		return "", fmt.Errorf("unknown otp type: %s", otpType)
	}
}

var ctx = context.Background()

func (r *otpRepository) SetOtp(otpType string, email string) (string, error) {
	otp, err := createOtp(otpType)
	if err != nil {
		return "", err
	}
	err = r.rdb.Set(ctx, otpType+":"+email, otp, time.Minute*15).Err()
	if err != nil {
		return "", err
	}
	return otp, nil
}

func (r *otpRepository) CheckOtp(otpType string, email, otp string) bool {
	val, err := r.rdb.Get(ctx, otpType+":"+email).Result()
	if err != nil {
		return false
	}
	return val != otp
}
