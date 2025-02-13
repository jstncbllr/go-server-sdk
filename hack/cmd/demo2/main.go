package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/launchdarkly/go-server-sdk/v7/hack/hacksdk"
)

func fetchSchemeData() (string, error) {
	resp, err := http.Get("http://localhost:8123/gonf")
	if err != nil {
		return "", fmt.Errorf("error fetching from proxy: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return string(body), nil
}

func main() {
	url := "http://localhost:8123/gonf" // Replace with your endpoint
	pollInterval := 5 * time.Second     // Poll every 5 seconds

	sdk, err := hacksdk.NewSDK()
	if err != nil {
		panic(err)
	}
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

	schemeData, err := fetchSchemeData()
	if err != nil {
		panic(err)
	}
	amount, err := sdk.Eval(schemeData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Congratulations! You owe %d dollars\n", amount)
}
