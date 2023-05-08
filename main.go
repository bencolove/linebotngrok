package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var apiKey string

func main() {
	fmt.Println("done")

	apiKey = os.Getenv("APIKEY")

	if apiKey == "" {
		// try from env file
		if key, err := readConfigFile(); err != nil {
			fmt.Fprintln(os.Stderr, "ApiKey not set")
			os.Exit(1)
		} else {
			apiKey = key
		}
	}

	startServer().Run()

}

func readConfigFile() (string, error) {
	// look for env file in the same folder
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return "", err
	}
	return viper.GetString("ApiKey"), nil
}
