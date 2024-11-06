// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"static-credential-provider/internel/utils"
)

func main() {

	// Read input from stdin
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}

	// Parse the input JSON string
	image, err := utils.GetRequestImage(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		os.Exit(1)
	}
	// Read values from environment variables or conf file
	username := os.Getenv("KSCP_REGISTRY_USERNAME")
	password := os.Getenv("KSCP_REGISTRY_PASSWORD")
	cacheType := os.Getenv("KSCP_CACHE_TYPE")
	cacheDuration := os.Getenv("KSCP_CACHE_DURATION")
	if len(os.Args) > 2 {
		arg1, val1 := os.Args[1], os.Args[2]
		if arg1 == "--config" {
			config, err := utils.GetConfig(val1)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading config file:", err)
				os.Exit(1)
			}
			username = config.Username
			password = config.Password
			cacheType = config.CacheType
			cacheDuration = config.CacheDuration
		}
	}

	if username == "" || password == "" {
		fmt.Fprintln(os.Stderr, "Error reading username or password from environment variables")
		os.Exit(1)
	}
	// Defaults to image cache type

	resp := utils.CreateImageRequestResponse(image, username, password, cacheType, cacheDuration)
	fmt.Println(resp)
}
