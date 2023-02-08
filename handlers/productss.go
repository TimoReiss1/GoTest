package handlers

import (
	"log"
	"main/work/data"
	"net/http"
	"regexp"
	"strconv"
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
	//handle an update using restful based approach, expect ID in URI
	if r.Method == http.MethodPut {
		p.l.Println("PUT Request:")
		r2 := regexp.MustCompile(`/([0-9]+)`)
		g := r2.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "Invalid ID", http.StatusBadRequest)
			return
		}
		p.l.Println("got id", id)

		p.updateProduct(id, rw, r)
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

	if err == nil {
		p.l.Printf("Prod: %#v", prod)
		data.AddProduct(prod)
	}

}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Request for: Products")
	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	if err == nil {
		p.l.Printf("Prod: %#v", prod)
		err = data.UpdateProduct(id, prod)
		if err == data.ErrorProductNotFound {
			http.Error(rw, "Product not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(rw, "Product not found", http.StatusInternalServerError)
			return
		}
	}
}

//Example Data for POST: curl -v localhost:9000 -d "{\"id\": 1, \"name\": \"tea\", \"description\": \"a nice cup of tea\"}"
