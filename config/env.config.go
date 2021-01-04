package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
	// var env map[string]string
	var config EnvConfig

	if err := godotenv.Load(); err != nil {
		fmt.Printf("[Error]: could not load env file. %s", err)
	}

	config.Port = os.Getenv("PORT")
	config.DBName = os.Getenv("DB_NAME")
	config.DBUri = os.Getenv("DB_URI")
	config.ClientHostname = os.Getenv("CLIENT_HOSTNAME")
	config.JwtSecret = os.Getenv("JWT_SECRET")
	config.AwsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	config.AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	config.AwsRegion = os.Getenv("AWS_REGION")
	config.AwsS3Bucket = os.Getenv("AWS_S3_BUCKET")

	return &config
}
