package tests

import (
	"net/http"
	"strings"
	"testing"
)

type reqRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Otp      string `json:"otp"`
}

func TestAuthRegisterBadRequestCases(t *testing.T) {
	baseReq := reqRegister{
		Username: "validUser",
		Password: "ValidPass123!",
		Email:    "test@example.com",
		Otp:      "123456",
	}
	cases := []struct {
		name    string
		modify  func(*reqRegister)
		message string
	}{
		{
			name:    "InvalidEmail",
			modify:  func(r *reqRegister) { r.Email = "not-an-email" },
			message: "邮箱必须是有效的邮箱地址",
		},
		{
			name:    "ShortUsername",
			modify:  func(r *reqRegister) { r.Username = "a" },
			message: "用户名至少需要2个字符",
		},
		{
			name:    "ShortUsernameUnicode",
			modify:  func(r *reqRegister) { r.Username = "你" },
			message: "用户名至少需要2个字符",
		},
		{
			name:    "LongUsername",
			modify:  func(r *reqRegister) { r.Username = strings.Repeat("a", 17) },
			message: "用户名不能超过16个字符",
		},
		{
			name:    "InvalidUsername",
			modify:  func(r *reqRegister) { r.Username = " User" },
			message: "用户名前后不能有空格",
		},
		{
			name:    "InvalidUsername",
			modify:  func(r *reqRegister) { r.Username = "U\nser" },
			message: "用户名只能包含可打印字符",
		},
		{
			name:    "ShortVerifyCode",
			modify:  func(r *reqRegister) { r.Otp = "123" },
			message: "验证码长度必须为6位",
		},
		{
			name:    "ShortPassword",
			modify:  func(r *reqRegister) { r.Password = "short" },
			message: "密码至少需要8个字符",
		},
		{
			name:    "LongPassword",
			modify:  func(r *reqRegister) { r.Password = strings.Repeat("a", 101) },
			message: "密码不能超过100个字符",
		},
		{
			name:    "InvalidVerifyOtp",
			modify:  func(r *reqRegister) { r.Otp = "abcdef" },
			message: "验证码必须是数字",
		},
		{
			name:    "ShortVerifyOtp",
			modify:  func(r *reqRegister) { r.Otp = "123" },
			message: "验证码长度必须为6位",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := baseReq
			tc.modify(&req)
			SendRequestAndExpectError(
				t, http.MethodPost, "/api/v1/auth/register", req,
				http.StatusBadRequest, tc.message,
			)
		})
	}
}
