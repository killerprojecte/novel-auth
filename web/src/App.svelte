<script lang="ts">
  import { Toaster } from "svelte-french-toast";

  let type: string = $state("登录");
  let app = "webx";
</script>

<main>
  <Toaster />
  <div
    class="fixed inset-0 overflow-y-auto
      bg-[url('https://books.fishhawk.top/assets/banner-BtpB_r33.webp')]"
  >
    <div class="absolute inset-0 -z-10 bg-black/80"></div>

    <div
      class="mx-auto flex h-full w-full flex-col gap-4
        bg-white px-8 pt-[10vh] pb-8 text-center
        sm:mt-[10vh] sm:h-auto sm:w-md sm:rounded-2xl sm:p-8"
    >
      <img
        class="mx-auto aspect-square w-1/2 max-w-[200px] select-none"
        src="https://n.novelia.cc/files-extra/girl.6e4fe22c238737fd028247f8f0cfd4ee.webp"
        alt="Girl"
      />

      <Tabs tabs={["登录", "注册"]} bind:activeTab={type} />

      {#if type === "登录"}
        <FormLogin {app} openResetPasswordForm={() => (type = "重置密码")} />
      {:else if type === "注册"}
        <FormRegister {app} />
      {:else}
        <FormResetPassword openLoginForm={() => (type = "登录")} />
      {/if}
    </div>
  </div>
</main>
