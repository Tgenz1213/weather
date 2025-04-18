// @ts-check

import { resolve } from "node:path"

import eslint from "@eslint/js"
import stylistic from "@stylistic/eslint-plugin"
import eslintConfigPrettier from "eslint-config-prettier"
import pathAlias from "eslint-plugin-path-alias"
import simpleImportSort from "eslint-plugin-simple-import-sort"
import tseslint from "typescript-eslint"

export default tseslint.config(
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  stylistic.configs.recommended,
  eslintConfigPrettier,
  {
    plugins: {
      "path-alias": pathAlias,
      "simple-import-sort": simpleImportSort,
      "@stylistic": stylistic,
    },
    rules: {
      "path-alias/no-relative": [
        "error",
        {
          paths: {
            "@": resolve(import.meta.dirname, "./src"),
            "@assets": resolve(import.meta.dirname, "./src/assets"),
            "@components": resolve(import.meta.dirname, "./src/components"),
            "@hooks": resolve(import.meta.dirname, "./src/hooks"),
            "@routes": resolve(import.meta.dirname, "./src/routes"),
            "@services": resolve(import.meta.dirname, "./src/services"),
          },
          exceptions: ["*.css"],
        },
      ],
      "simple-import-sort/imports": "error",
      "simple-import-sort/exports": "error",
      "@stylistic/brace-style": ["error", "1tbs"],
    },
  },
  { ignores: ["**/dist/**", "**/node_modules/**"] },
)
