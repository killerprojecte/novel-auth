<script lang="ts">
  import { toast } from "svelte-sonner";
  import { Api, onLoginSuccess } from "../data/api";

  let { app, openResetPasswordForm } = $props();

  let username = $state("");
  let password = $state("");

  let loading = $state(false);
  function login(event: MouseEvent) {
    event.preventDefault();

    if (Api.login.isPending) return;
    loading = true;
    Api.login(app, username, password)
      .then(() => {
        loading = false;
        onLoginSuccess();
      })
      .catch((error) => {
        loading = false;
        toast.error(`登录失败: ${error}`);
      });
  }
</script>

<form class="flex w-auto flex-col gap-2" novalidate>
  <FormItem>
    <Input placeholder="用户名/邮箱" bind:value={username} />
  </FormItem>
  <FormItem>
    <Input type="password" placeholder="密码" bind:value={password} />
  </FormItem>

  <button
    class="text-primary cursor-pointer text-right text-sm font-bold"
    onclick={openResetPasswordForm}>忘记密码？</button
  >

  <Button {loading} text="登录" onclick={login} />
</form>
