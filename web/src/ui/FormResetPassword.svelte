<script lang="ts">
  import toast from "svelte-french-toast";
  import { Api, onLoginSuccess } from "../data/api";
  import { Validator } from "./util";
  import OtpButton from "./OtpButton.svelte";

  let password = $state("");
  let email = $state("");
  let otp = $state("");

  let loading = $state(false);
  function resetPassword(event: MouseEvent) {
    event.preventDefault();

    if (Api.resetPassword.isPending) return;
    loading = true;
    Api.resetPassword(email, password, otp)
      .then(() => onLoginSuccess())
      .catch((error) => {
        loading = false;
        toast.error(`重置密码失败: ${error}`);
      });
  }
</script>

<form class="flex w-auto flex-col gap-2" novalidate>
  <FormItem rules={Validator.validateEmail}>
    <Input placeholder="邮箱" bind:value={email} />
  </FormItem>

  <FormItem rules={Validator.validateOtp}>
    <Input round="left" placeholder="邮箱验证码" bind:value={otp} />
    <OtpButton {email} type="reset_password" round="right" class="flex-1/2" />
  </FormItem>

  <FormItem rules={Validator.validatePassword}>
    <Input placeholder="新密码" bind:value={password} />
  </FormItem>

  <p class="mt-1 text-left text-xs text-[#8d8d8d] select-none">
    * 收不到验证邮件的话，记得看垃圾箱
  </p>

  <Button text="重置密码" {loading} onclick={resetPassword} />
</form>
