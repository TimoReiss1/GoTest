package main

import (
	"context"
	"log"
	"main/work/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Hello world, the web server
	/*
		http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			log.Println("Hello World!")
			d, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(rw, "Oops", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(rw, "Hello %s\n", d)
		})

		helloHandler := func(w http.ResponseWriter, req *http.Request) {
			io.WriteString(w, "Hello, world!\n")
		}
		http.HandleFunc("/hello", helloHandler)

		goodbyeHandler := func(w http.ResponseWriter, req *http.Request) {
			io.WriteString(w, "Bye World!\n")
		}
		http.HandleFunc("/bye", goodbyeHandler)
	*/
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	//hh := handlers.NewHello(l)
	//gh := handlers.NewBye(l)
	//sm.Handle("/", hh)
	//sm.Handle("/goodbye", gh)
	ph := handlers.NewProducts(l)
	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := &http.Server{
		Addr:         ":9000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	log.Println("WebServer started successfully")
	//http.ListenAndServe(":9000", sm)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, gracefully shuttingdown.\n Signal Type: ", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
