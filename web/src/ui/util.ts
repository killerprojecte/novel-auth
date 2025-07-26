export namespace Validator {
  export function validateUsername(username: string) {
    if (!username) return "用户名不能为空";
    if (username.length < 3) return "用户名至少 3 个字符";
    return true;
  }
  export function validatePassword(password: string) {
    if (!password) return "密码不能为空";
    if (password.length < 6) return "密码至少 6 个字符";
    return true;
  }
  export function validateEmail(email: string) {
    if (!email) return "邮箱不能为空";
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) return "请输入有效的邮箱地址";
    return true;
  }
  export function validateOtp(otp: string) {
    if (!otp) return "邮箱验证码不能为空";
    if (otp.length !== 6) return "邮箱验证码必须是 6 位数字";
    return true;
  }
}
