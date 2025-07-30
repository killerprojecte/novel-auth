<script lang="ts">
  import { toast } from "svelte-sonner";
  import { Api } from "../data/api";

  let { type, email, class: className, ...rest } = $props();

  let countdown = $state(0);

  function startCountdown() {
    countdown = 60;
    const interval = setInterval(() => {
      if (countdown > 0) {
        countdown -= 1;
      } else {
        clearInterval(interval);
      }
    }, 1000);
  }

  function requestOtp(event: MouseEvent) {
    event.preventDefault();

    if (Api.requestOtp.isPending || countdown > 0) return;
    Api.requestOtp(email, type)
      .then(() => {
        toast.success("验证码已发送到您的邮箱");
        startCountdown();
      })
      .catch((error) => {
        toast.error(`验证码发送失败: ${error}`);
        startCountdown();
      });
  }
</script>

<Button
  {...rest}
  disabled={Api.requestOtp.isPending || countdown > 0}
  text={countdown > 0 ? `${countdown}秒冷却` : "发送验证码"}
  round="right"
  onclick={requestOtp}
  class={className}
/>
