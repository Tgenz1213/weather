package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	env := flag.String("env", "", "Environment to run the application. Default value is 'example' for .env.[example].")

	flag.Parse()

	var envFile string
	if *env == "" {
		envFile = ".env"
	} else {
		envFile = ".env." + *env
	}

	currentWorkDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current work directory: %v", err)
	}

	envFilePath := filepath.Join(currentWorkDir, "..", "..", envFile)

	err = godotenv.Load(envFilePath)
	if err != nil {
		return fmt.Errorf("cannot find env file: %v", err)
	}

	return nil
}
