package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Bid string `json:"bid"`
}

func main() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("Error parsing response:", err)
	}

	// Write to file
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Dollar: %s", response.Bid)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}

	fmt.Printf("Exchange rate saved to cotacao.txt: %s\n", response.Bid)
}
