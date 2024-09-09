import { defineConfig } from "vitest/config";

// refer to: https://vitest.dev/config/
export default defineConfig({
  test: {
    // clearMocks: true,
    // globals: true,
    isolate: false,
    setupFiles: [`${__dirname}/vitest.setup.ts`],
  },
});
