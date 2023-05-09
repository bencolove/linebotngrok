package main

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func GetEnvString(key string) (string, error) {

	if val := os.Getenv(strings.ToUpper(key)); val != "" {
		return val, nil
	} else {
		return readConfigFile(key)
	}
}

var isConfigRead = false

func readConfigFile(key string) (string, error) {
	if isConfigRead == false {
		// look for env file in the same folder
		viper.SetConfigName("env")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			return "", err
		}
		isConfigRead = true
	}

	return viper.GetString(key), nil
}
