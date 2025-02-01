import swc from 'unplugin-swc';
import { defineConfig } from 'vitest/config';

// refer to: https://vitest.dev/config/
export default defineConfig({
  plugins: [swc.vite()],
  test: {
    clearMocks: true,
    // globals: true,
    isolate: false,
    setupFiles: [`${__dirname}/vitest.setup.mts`],
    coverage: {
      provider: "v8",
      include: ["src/**/*"],
      exclude: [
        "**/dist/**", // exclude any dist folder under src
        // and the vitest coverage.exclude defaults
        "coverage/**",
        "dist/**",
        "**\/[.]**",
        "packages/*\/test?(s)/**",
        "**\/*.d.ts",
        "**\/virtual:*",
        "**\/__x00__*",
        "**\/\x00*",
        "cypress/**",
        "test?(s)/**",
        "test?(-*).?(c|m)[jt]s?(x)",
        "**\/*{.,-}{test,spec}?(-d).?(c|m)[jt]s?(x)",
        "**\/__tests__/**",
        "**\/{karma,rollup,webpack,vite,vitest,jest,ava,babel,nyc,cypress,tsup,build}.config.*",
        "**\/vitest.{workspace,projects}.[jt]s?(on)",
        "**\/.{eslint,mocha,prettier}rc.{?(c|m)js,yml}",
      ],
    },
  }
});
