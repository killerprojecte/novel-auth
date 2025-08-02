import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import Components from "unplugin-svelte-components/vite";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  server: {
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "*",
      "Access-Control-Allow-Headers": "Content-Type",
    },
  },
  plugins: [
    tailwindcss(),
    svelte(),
    Components({
      dts: "src/components.d.ts",
      dirs: ["src/components", "src/ui"],
    }),
  ],
});
