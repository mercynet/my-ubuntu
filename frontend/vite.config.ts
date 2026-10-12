import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Everything is bundled locally (no CDN); Go embeds dist/ with go:embed.
export default defineConfig({
  plugins: [react()],
  base: "./",
  build: { outDir: "dist", emptyOutDir: true, assetsInlineLimit: 0, chunkSizeWarningLimit: 1500 },
});
