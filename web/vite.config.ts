import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import Components from "unplugin-svelte-components/vite";
import { defineConfig } from "vite";
import { viteSingleFile } from "vite-plugin-singlefile";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss(),
    svelte(),
    viteSingleFile(),
    Components({
      dts: "src/components.d.ts",
      dirs: ["src/components", "src/ui"],
    }),
  ],
});
