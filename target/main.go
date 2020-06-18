package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	counter := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/root", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// some important business logic
		counter++
		timestamp := r.Header.Get("h1")
		response := fmt.Sprintf("timestamp %s : count %d", timestamp, counter)

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, response)
		return
	})

	// this is the target server
	srv := http.Server{Addr: ":10001", Handler: mux}
	fmt.Println("starting target server")

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
