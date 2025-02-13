package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	url := "http://localhost:8123/gonf" // Replace with your endpoint
	pollInterval := 5 * time.Second     // Poll every 5 seconds

	go func() {
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
	}()

	fmt.Println("Time to pay your taxes! Please enter your income 💰")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	income := scanner.Text()
	strconv.ParseFloat(income, 64)

  amount := sdk.Evaluate()
}
