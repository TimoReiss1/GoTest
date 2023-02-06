package main

import (
	"log"
	"main/handlers"
	"net/http"
	"os"
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
	hh := handlers.NewHello(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)

	log.Println("Listing for requests at http://localhost:9000/hello")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
