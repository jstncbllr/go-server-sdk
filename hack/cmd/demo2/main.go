package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/markkurossi/scheme"

	"github.com/launchdarkly/go-server-sdk/v7/hack/hacksdk"
)

func fetchSchemeData() (string, error) {
	resp, err := http.Get("http://localhost:8123/other")
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
	url := "http://localhost:8123/other" // Replace with your endpoint
	pollInterval := 5 * time.Second      // Poll every 5 seconds

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

			_, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Printf("Error reading response: %v\n", err)
				time.Sleep(pollInterval)
				continue
			}

			// fmt.Printf("Response: %s\n", string(body))
			time.Sleep(pollInterval)
		}
	}()

	fmt.Println("Time to pay your taxes! Please enter your income 💰")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	incomeStr := scanner.Text()
	income, err := strconv.ParseFloat(incomeStr, 64)
	if err != nil {
		panic(err)
	}

	// 	_, err = sdk.Eval(`
	// 		(define (assoc key lst)
	// 			(cond
	// 				((null? lst) #f)
	// 				((equal? key (car (car lst))) (car lst))
	// 				(else (assoc key (cdr lst)))))
	//
	// 		(define (assoc-cdr key lst)
	// 			(cdr (assoc key lst)))
	// 	`)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	schemeData, err := fetchSchemeData()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Scheme data: %s\n\n", schemeData)
	_, err = sdk.Eval(schemeData)
	if err != nil {
		panic(err)
	}
	percent, err := sdk.Eval(fmt.Sprintf(`(evaluate %f)`, income))
	if err != nil {
		panic(err)
	}
	percentStr := scheme.ToString(percent)
	asNum, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Congratulations! You owe $%.2f dollars\n", asNum)
}
