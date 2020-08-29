package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// EnvConfig -
type EnvConfig struct {
	DBName string `mapstructure:"DB_NAME"`
	Port   int    `mapstructure:"PORT"`
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

	fmt.Printf("[DEBUG]: config - %+v", env)

	return env
}
