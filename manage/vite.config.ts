import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

const srcAlias = decodeURIComponent(
  new URL("./src", import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"),
);

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": srcAlias,
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8002",
        changeOrigin: true,
      },
    },
  },
});
