import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";

const srcAlias = decodeURIComponent(
    new URL("./src", import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"),
);

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, ".", "");
  const devApiProxyTarget =
      env.VITE_DEV_API_PROXY_TARGET || "http://127.0.0.1:8002";

  return {
    // 关键：后台部署在 /admin/ 下
    base: "/admin/",

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
          target: devApiProxyTarget,
          changeOrigin: true,
        },
      },
    },
  };
});