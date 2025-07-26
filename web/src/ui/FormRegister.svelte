<script lang="ts">
  import { Api, redirectAfterLogin } from "../data/api";
  import { Validator } from "./util";

  let { app } = $props();

  let username = $state("");
  let password = $state("");
  let email = $state("");
  let otp = $state("");

  let loading = $state(false);
  function register(event: MouseEvent) {
    event.preventDefault();

    if (Api.register.isPending) return;
    loading = true;
    Api.register(app, email, username, password, otp)
      .then(() => redirectAfterLogin())
      .catch((error) => {
        loading = false;
        alert(`注册失败: ${error}`);
      });
  }

  function requestOtp(event: MouseEvent) {
    event.preventDefault();

    if (Api.requestOtp.isPending) return;
    Api.requestOtp(email)
      .then(() => {
        alert("验证码已发送到您的邮箱");
      })
      .catch((error) => {
        alert(`发送验证码失败: ${error}`);
      });
  }
</script>

<form class="flex w-auto flex-col gap-2" novalidate>
  <FormItem rules={Validator.validateUsername}>
    <Input placeholder="用户名" bind:value={username} />
  </FormItem>

  <FormItem rules={Validator.validatePassword}>
    <Input placeholder="密码" bind:value={password} />
  </FormItem>

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

  <p class="mt-1 text-left text-xs text-[#8d8d8d] select-none">
    * 收不到验证邮件的话，记得看垃圾箱
  </p>

  <Button text="注册" {loading} onclick={register} />
</form>
