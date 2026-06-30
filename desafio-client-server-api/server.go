package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExchangeRate struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS exchange_rates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		// Context with timeout for the entire request
		ctx := r.Context()
		
		// Create context with timeout for API call
		apiCtx, apiCancel := context.WithTimeout(ctx, 200*time.Millisecond)
		defer apiCancel()

		// Make request to external API
		req, err := http.NewRequestWithContext(apiCtx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to get exchange rate: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var exchangeRate ExchangeRate
		if err := json.NewDecoder(resp.Body).Decode(&exchangeRate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Context with timeout for database operation
		dbCtx, dbCancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer dbCancel()

		// Save to database
		_, err = db.ExecContext(dbCtx,
			"INSERT INTO exchange_rates (bid) VALUES (?)",
			exchangeRate.USDBRL.Bid,
		)
		if err != nil {
			log.Printf("Failed to save to database: %v", err)
			http.Error(w, "Database operation timeout", http.StatusInternalServerError)
			return
		}

		// Return only the bid value as JSON
		response := struct {
			Bid string `json:"bid"`
		}{
			Bid: exchangeRate.USDBRL.Bid,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
