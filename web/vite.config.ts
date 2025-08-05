import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { defineConfig } from 'vite';
import { createHtmlPlugin } from 'vite-plugin-html';

// https://vite.dev/config/
export default defineConfig({
  build: {
    target: ['es2015'],
  },
  server: {
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': '*',
      'Access-Control-Allow-Headers': 'Content-Type',
    },
  },
  plugins: [
    tailwindcss(),
    vue(),
    createHtmlPlugin({
      minify: { minifyJS: true },
    }),
    AutoImport({
      dts: './src/auto-imports.d.ts',
      imports: ['vue', 'vue-router', 'pinia'],
    }),
    Components({
      dts: 'src/components.d.ts',
      dirs: ['src/components', 'src/ui'],
    }),
  ],
});
