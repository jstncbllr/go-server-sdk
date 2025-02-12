package cmd

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:8080/your-endpoint" // Replace with your endpoint
	pollInterval := 5 * time.Second              // Poll every 5 seconds

	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			time.Sleep(pollInterval)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			time.Sleep(pollInterval)
			continue
		}

		fmt.Printf("Response: %s\n", string(body))
		time.Sleep(pollInterval)
	}
}
