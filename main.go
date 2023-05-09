package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("done")

	if apiKey, err := GetEnvString("ApiKey"); apiKey == "" || err != nil {
		if err != nil {
			fmt.Fprintln(os.Stderr, "ApiKey not set")
		}
		if apiKey == "" {
			fmt.Fprintln(os.Stderr, "ApiKey not set")
		}
		os.Exit(1)
	}

	if server, err := startServer(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	} else {
		server.Run()
	}

}
