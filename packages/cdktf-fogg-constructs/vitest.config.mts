import swc from "unplugin-swc";
import { defineConfig } from "vitest/config";

// refer to: https://vitest.dev/config/
export default defineConfig({
  plugins: [swc.vite()],
  test: {
    clearMocks: true,
    // globals: true,
    isolate: false,
    setupFiles: [`${__dirname}/vitest.setup.mts`],
  },
});
