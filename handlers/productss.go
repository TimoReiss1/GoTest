package handlers

import (
	"log"
	"main/work/data"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

/*
func (p *Products) ServeHTTP(rw http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
*/

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//handle a post request
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	//handle an update using restful based approach
	if r.Method == http.MethodPut {
		p.updateProduct(rw, r)
		return
	}
	//catch rest
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, hr *http.Request) {
	p.l.Println("Handle GET Request for: Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Request for: Products")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Request for: Products")
}

//Example Data for POST: curl -v localhost:9000 -d "{\"id\": 1, \"name\": \"tea\", \"description\": \"a nice cup of tea\"}"
