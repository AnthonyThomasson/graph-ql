import vue from "@vitejs/plugin-vue";
import { defineConfig } from "vite";
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    manifest: true,
    rollupOptions: {
      input: "src/main.ts",
    },
  },
});
