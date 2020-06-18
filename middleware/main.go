package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/root", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// initialize client
		client := http.Client{Timeout: time.Second * 30}

		// create a new request
		request, err := http.NewRequest("GET", "http://localhost:10001/root", nil)
		if err != nil {
			log.Println("error creating request", err)
		}
		request.Header.Set("h1", time.Now().Format("02-Jan-2006 15:04:05"))

		// send a request
		res, err := client.Do(request)
		if err != nil {
			log.Println(err)
			_, _ = io.WriteString(w, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// read response
		defer func() { _ = res.Body.Close() }()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil || res.StatusCode != http.StatusOK {
			log.Println(err)
			_, _ = io.WriteString(w, "error response from server")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// middleware responds with target server's response
		resp := string(b)
		defer log.Println(resp)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, resp)
		return
	})

	log.Println("starting standard middleware")
	srv := http.Server{Addr: ":10000", Handler: mux}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
