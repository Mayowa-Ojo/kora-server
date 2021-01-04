package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
)

// EnvConfig -
type EnvConfig struct {
	Port               string `mapstructure:"PORT"`
	DBName             string `mapstructure:"DB_NAME"`
	DBUri              string `mapstructure:"DB_URI"`
	ClientHostname     string `mapstructure:"CLIENT_HOSTNAME"`
	JwtSecret          string `mapstructure:"JWT_SECRET"`
	AwsAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion          string `mapstructure:"AWS_REGION"`
	AwsS3Bucket        string `mapstructure:"AWS_S3_BUCKET"`
}

// NewEnvConfig -
func NewEnvConfig() *EnvConfig {
	var env map[string]string
	var config EnvConfig

	if err := godotenv.Load(); err != nil {
		fmt.Printf("[Error]: could not load env file. %s", err)
	}

	env, err := godotenv.Read()
	if err != nil {
		fmt.Printf("[Error]: could not read env file. %s", err)
	}

	if err = mapstructure.Decode(env, &config); err != nil {
		fmt.Printf("[Error]: error decoding map structure. %s", err)
	}

	return &config
}
