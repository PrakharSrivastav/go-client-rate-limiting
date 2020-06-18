package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// pool collects all the incoming http requests and stores them as jobs
type pool struct {
	jobs   chan *http.Request
	client *http.Client
	done   chan struct{}
}

func setup() *pool {
	// estimate the maximum size of requests you can get per unit of time
	// so that we can respond to the called ASAP
	ch := make(chan *http.Request, 5000)
	cl := &http.Client{Timeout: 5 * time.Second}
	return &pool{
		jobs:   ch,
		client: cl,
		done:   make(chan struct{}),
	}
}

// this is to force the client to block every "rateLimit" Seconds
const rateLimit = time.Second

// start this processor to send requests to target server.
// It is blocking in nature so make sure to run it in a separate goroutine
func (p *pool) start() {
	// throttle is a handle that will unblock evert "rateLimit" seconds
	throttle := time.Tick(rateLimit)
	for {
		// block right away
		<-throttle

		select {
		case job := <-p.jobs:
			p.doJob(job)
		case <-p.done:
			// graceful shutdown
			// cleanup resources
			log.Println("stopping")

			// this will exit the goroutine
			return
		}
	}

}

// perform your job (call target server) here
func (p *pool) doJob(job *http.Request) {
	res, err := p.client.Do(job)
	if err != nil {
		log.Println(err)
	}
	defer func() { _ = res.Body.Close() }()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	log.Printf("completed job : %s \n", string(b))
}

// add requests to the pool
func (p *pool) submit(req *http.Request) {
	p.jobs <- req
}

// use this to shutdown the pool
func (p *pool) close() {
	close(p.done)
}

func main() {
	p := setup()
	go p.start()

	mux := http.NewServeMux()
	mux.HandleFunc("/root", func(w http.ResponseWriter, r *http.Request) {

		request, err := http.NewRequest("GET", "http://localhost:10001/root", nil)
		if err != nil {
			log.Println(err)
		}
		request.Header.Set("h1", time.Now().UTC().String())

		// submit the request to pool
		p.submit(request)

		// and immediately reply with http.StatusOK
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	})
	log.Println("starting rate limited middleware")
	srv := http.Server{Addr: ":10000", Handler: mux}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
