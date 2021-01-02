package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// EnvConfig -
type EnvConfig struct {
	DBName             string `mapstructure:"DB_NAME"`
	Port               int    `mapstructure:"PORT"`
	ClientHostname     string `mapstructure:"CLIENT_HOSTNAME"`
	JwtSecret          string `mapstructure:"JWT_SECRET"`
	AwsAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion          string `mapstructure:"AWS_REGION"`
	AwsS3Bucket        string `mapstructure:"AWS_S3_BUCKET"`
}

// NewEnvConfig -
func NewEnvConfig() *EnvConfig {
	var env *EnvConfig
	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("[Error]: could not read env file. %s", err)
	}

	if err := v.Unmarshal(&env); err != nil {
		fmt.Printf("[Error]: could not decode env variables. %s", err)
	}

	return env
}
