package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/launchdarkly/go-server-sdk/v7/hack/json2scheme"
)

func main() {
	// Create handler for /gonf endpoint
	http.HandleFunc("/gonf", func(w http.ResponseWriter, r *http.Request) {
		// Create a new HTTP client
		client := &http.Client{}

		// Create the request to LaunchDarkly
		req, err := http.NewRequest("GET", "https://ld-stg.launchdarkly.com/sdk/latest-all", nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
			return
		}

		// Add the Authorization header
		req.Header.Add("Authorization", "sdk-68035835-ce25-4570-9eed-b71f5e86b0b4")

		// Make the request
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error making request: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError)
			return
		}

		// Convert JSON to Scheme
		schemeData, err := json2scheme.JsonToScheme(body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error converting to scheme: %v", err), http.StatusInternalServerError)
			return
		}

		// Wrap the scheme data in a payload definition
		wrappedScheme := fmt.Sprintf("(define payload '%s)", schemeData)

		// Add the evaluate function definition
		evaluateFunc := `(define (evaluate ctx flag-key default)(assoc-cdr 'on (car (assoc-cdr flag-key (car (assoc-cdr 'flags payload))))))`
		finalScheme := wrappedScheme + "\n" + evaluateFunc
		// finalScheme := wrappedScheme

		// Set content type header
		w.Header().Set("Content-Type", "text/plain")

		// Write the scheme response
		fmt.Fprint(w, finalScheme)
	})

	http.HandleFunc("/other", func(w http.ResponseWriter, r *http.Request) {
		payload := `
		(define (evaluate ctx)
			(/ (* ctx (find-bracket ctx)) 100)
		)
		(define (find-bracket ctx)
			(define (find-bracket-helper brackets)
				(if (null? brackets)
					nil
					(if (<= (caaar brackets) ctx)
						(cdr (caar brackets))
						(find-bracket-helper (cdr brackets)))))
			(find-bracket-helper (assoc-cdr 'brackets payload))
			)
		(define payload
			'((brackets
				((20000 0)
				 (30000 10)
				 (40000 15)
				 (50000 20)
				 (100000 30))
			))
		)
		`
		fmt.Fprint(w, payload)
	})

	// Start the server on port 8123
	port := ":8123"
	fmt.Printf("Starting server on %s/gonf\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
