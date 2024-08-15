package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	env := flag.String("env", "example", "Environment to run the application. Default value is 'example' for .env.[example].")

	flag.Parse()

	envFile := ".env." + *env

	currentWorkDir, err := os.Getwd()

	if err != nil {
		return err
	}

	envFilePath := filepath.Join(currentWorkDir, "..", "..", envFile)

	err = godotenv.Load(envFilePath)

	if err != nil {
		return err
	}

	return nil
}
