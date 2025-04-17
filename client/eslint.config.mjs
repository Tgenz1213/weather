// @ts-check

import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import pathAlias from "eslint-plugin-path-alias";
import { resolve } from "node:path";

export default tseslint.config(
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  {
    plugins: {
      "path-alias": pathAlias,
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
          exceptions: ["*.css"]
        },
      ],
    },
  },
  { ignores: ["**/dist/**", "**/node_modules/**"] }
);
