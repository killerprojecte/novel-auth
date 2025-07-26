<script lang="ts">
  import { Api, redirectAfterLogin } from "../data/api";
  import { Validator } from "./util";

  let password = $state("");
  let email = $state("");
  let otp = $state("");

  let loading = $state(false);
  function resetPassword(event: MouseEvent) {
    event.preventDefault();

    if (Api.resetPassword.isPending) return;
    loading = true;
    Api.resetPassword(email, password, otp)
      .then(() => redirectAfterLogin())
      .catch((error) => {
        loading = false;
        alert(`重置密码失败: ${error}`);
      });
  }

  function requestOtp(event: MouseEvent) {
    event.preventDefault();

    if (Api.requestOtp.isPending) return;
    Api.requestOtp(email, "reset_password")
      .then(() => {
        alert("验证码已发送到您的邮箱");
      })
      .catch((error) => {
        alert(`发送验证码失败: ${error}`);
      });
  }
</script>

<form class="flex w-auto flex-col gap-2" novalidate>
  <FormItem rules={Validator.validateEmail}>
    <Input placeholder="邮箱" bind:value={email} />
  </FormItem>

  <FormItem rules={Validator.validateOtp}>
    <Input round="left" placeholder="邮箱验证码" bind:value={otp} />
    <Button
      text="发送验证码"
      round="right"
      onclick={requestOtp}
      class="flex-1/2"
    />
  </FormItem>

  <FormItem rules={Validator.validatePassword}>
    <Input placeholder="新密码" bind:value={password} />
  </FormItem>

  <p class="mt-1 text-left text-xs text-[#8d8d8d] select-none">
    * 收不到验证邮件的话，记得看垃圾箱
  </p>

  <Button text="重置密码" {loading} onclick={resetPassword} />
</form>
