package main

import (
	"fmt"
	"io"
	"net/http"
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

func printBoolAscii(value string) {
	if value == "(#t)" {
		fmt.Println(`
 _____  ____  _    _ _____ 
|_   _||  _ \| |  | |  ___|
  | |  | |_) | |  | | |__  
  | |  |  _ <| |  | |  __| 
  | |  | |_) | |__| | |___ 
  |_|  |____/ \____/|_____|
`)
	} else if value == "(#f)" {
		fmt.Println(`
 _____ ___  _     ____  _____ 
|  ___/ _ \| |   / ___|| ____|
| |_ | | | | |   \___ \|  _|  
|  _|| |_| | |___ ___) | |___ 
|_|   \___/|_____|____/|_____|
`)
	} else {
		panic(value)
	}
}

func main() {
	pollInterval := 5 * time.Second // Poll every 5 seconds
	for {

		sdk, err := hacksdk.NewSDK()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Loading helpers\n")

		// Define our helper functions first
		_, err = sdk.Eval(`
		(define (assoc key lst)
			(cond
				((null? lst) #f)
				((equal? key (car (car lst))) (car lst))
				(else (assoc key (cdr lst)))))

		(define (assoc-cdr key lst)
			(cdr (assoc key lst)))
	`)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Calling gonf proxy\n")

		// Fetch scheme data from proxy
		schemeData, err := fetchSchemeData()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Loading gonf response\n")

		// Evaluate the fetched scheme data
		_, err = sdk.Eval(schemeData)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Evaluating flag 'foo'\n")

		// Evaluate a specific flag
		result, err := sdk.Eval(`(evaluate "" 'foo #f)`)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Flag value at %v: ", time.Now().Format(time.RFC3339))
		printBoolAscii(result.Scheme())
		time.Sleep(pollInterval)
	}
}
