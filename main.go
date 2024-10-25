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
		fmt.Println("Error reading input:", err)
		return
	}

	// Parse the input JSON string
	image, err := utils.GetRequestImage(input)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	// Read values from environment variables or conf file
	username := os.Getenv("SCP_REGISTRY_USERNAME")
	password := os.Getenv("SCP_REGISTRY_PASSWORD")
	cacheType := os.Getenv("CACHE_TYPE")
	if len(os.Args) > 2 {
		arg1, val1 := os.Args[1], os.Args[2]
		if arg1 == "--config" {
			config, err := utils.GetConfig(val1)
			if err != nil {
				fmt.Println("Error reading config file:", err)
				return
			}
			username = config.Username
			password = config.Password
			cacheType = config.CacheType
		}
	}

	if username == "" || password == "" {
		fmt.Println("Error reading username or password from environment variables")
		return
	}
	// Defaults to image cache type

	resp := utils.CreateImageRequestResponse(image, username, password, cacheType)
	fmt.Println(resp)
}
