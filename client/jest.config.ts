import type { Config } from "jest";
import { configDotenv } from "dotenv";

configDotenv({ path: "./.env.development" });

const config: Config = {
  extensionsToTreatAsEsm: [".ts"],
  verbose: true,
  testEnvironment: "jest-environment-jsdom",
  testPathIgnorePatterns: ["/node_modules/", "/dist/"],
  transform: {
    "^.+.tsx?$": "babel-jest",
  },
  moduleNameMapper: {
    "\\.(css|less)$": "identity-obj-proxy",
  },
};

export default config;
