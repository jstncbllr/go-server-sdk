package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/launchdarkly/go-server-sdk/v7/hack/json2scheme"
)

func main() {
	// Create a new HTTP client
	client := &http.Client{}

	// Create the request
	req, err := http.NewRequest("GET", "https://ld-stg.launchdarkly.com/sdk/latest-all", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Add the Authorization header
	req.Header.Add("Authorization", "sdk-68035835-ce25-4570-9eed-b71f5e86b0b4")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Print the response
	fmt.Println(string(body))

	scheme, err := json2scheme.JsonToScheme(body)
	if err == nil {
		println(err)
	}
	fmt.Println(scheme)

}
