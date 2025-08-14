import postcssCascadeLayers from '@csstools/postcss-cascade-layers';
import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { defineConfig, type UserConfig } from 'vite';
import { createHtmlPlugin } from 'vite-plugin-html';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const config: UserConfig = {
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
    css: {},
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
  };

  // postcss-cascade-layers 在开发模式下会导致样式加载异常，因此仅在生产模式下启用
  if (mode === 'production') {
    config.css!.postcss = {
      plugins: [postcssCascadeLayers()],
    };
  }
  return config;
});
