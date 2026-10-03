import { copyFile } from "node:fs/promises";
import { defineConfig, lazyPlugins } from "vite-plus";
import react from "@vitejs/plugin-react";

const githubPagesFallback = () => ({
  name: "github-pages-fallback",
  closeBundle: async () => {
    await copyFile("dist/index.html", "dist/404.html");
  },
});

export default defineConfig({
  base: "/mygo-breeze-ui/",
  plugins: lazyPlugins(() => [react(), githubPagesFallback()]),
  server: { port: 4173 },
});
