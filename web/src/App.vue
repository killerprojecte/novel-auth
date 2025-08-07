<script setup lang="ts">
import 'vue-sonner/style.css';
import { Toaster } from 'vue-sonner';

const type = ref('登录');

function getOsTheme() {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light';
}

const query = new URLSearchParams(window.location.search);
const app = query.get('app') || 'auth';
const theme = query.get('theme') || getOsTheme();

const root = window.document.documentElement;
if (theme === 'dark') {
  root.classList.add('dark');
} else {
  root.classList.remove('dark');
}
</script>

<template>
  <Toaster position="top-center" />
  <div
    class="fixed inset-0 overflow-y-auto bg-[url('https://books.fishhawk.top/assets/banner-BtpB_r33.webp')]"
  >
    <div class="absolute inset-0 -z-10 bg-black/80"></div>

    <div
      class="bg-background mx-auto flex h-full w-full flex-col gap-4 px-8 pt-[10vh] pb-8 text-center sm:mt-[10vh] sm:h-auto sm:w-md sm:rounded-2xl sm:p-8"
    >
      <img
        class="mx-auto aspect-square w-1/2 max-w-[200px] select-none"
        src="https://n.novelia.cc/files-extra/girl.6e4fe22c238737fd028247f8f0cfd4ee.webp"
        alt=""
      />

      <Tabs :tabs="['登录', '注册']" v-model="type" />

      <FormLogin
        v-if="type === '登录'"
        :app="app"
        @openResetPasswordForm="type = '重置密码'"
      />

      <FormRegister v-else-if="type === '注册'" :app="app" />

      <FormResetPassword
        v-else-if="type === '重置密码'"
        @openLoginForm="type = '登录'"
      />
    </div>
  </div>
</template>
