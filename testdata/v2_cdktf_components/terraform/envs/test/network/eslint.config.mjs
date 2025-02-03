import js from "@eslint/js";
import eslintConfigPrettier from "eslint-config-prettier";
import turboPlugin from "eslint-plugin-turbo";
import tseslint from "typescript-eslint";
import onlyWarn from "eslint-plugin-only-warn";
import globals from "globals";

/**
 * @type {import("eslint").Linter.Config}
 * */
export default [
  js.configs.recommended,
  eslintConfigPrettier,
  ...tseslint.configs.recommended,
  {
    plugins: {
      turbo: turboPlugin,
    },
    rules: {
      "no-console": "off",
      "no-constant-condition": "warn",
      "turbo/no-undeclared-env-vars": "off",
      "@typescript-eslint/no-inferrable-types": "off",
      "@typescript-eslint/ban-types": "off",
      "@typescript-eslint/no-explicit-any": "off",
    },
  },
  {
    plugins: {
      onlyWarn,
    },
  },
  {
    rules: {
      "no-undef": "off",
      "no-unused-vars": "off",
    },
    languageOptions: {
      parserOptions: {
        project: "tsconfig.json",
      },
      globals: {
        ...globals.node,
      },
    },
  },
  {
    ignores: [
      // Ignore dotfiles
      ".*.?(c)js",
      "*.config*.?(c|m)js",
      ".*.ts",
      "*.config*.ts",
      "*.d.ts",
      "**/dist/**",
      ".git/**",
      "**/node_modules/**",
      "coverage/**",
      "**/.pnpm/**",
    ],
  },
];
