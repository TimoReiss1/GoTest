package handlers

import (
	"database/sql"
	"log"
	"main/work/data"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Products struct {
	l *log.Logger
}

var apple *sql.DB

func NewProducts(l *log.Logger, db *sql.DB) *Products {
	apple = db
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
		//p.getProducts(rw, r)
		return
	}

	//handle a post request
	if r.Method == http.MethodPost {
		//p.AddProduct(rw, r)
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
		p1, _, _ := data.FindProductById(id - 1)
		p.l.Printf("Updated Product by ID: %#v", p1)

		//p.UpdateProduct(rw, r)
		return
	}
	//catch rest
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) GetProducts(rw http.ResponseWriter, hr *http.Request, _ httprouter.Params) {
	// Log a message indicating that a GET request for products was received
	p.l.Println("Handle GET Request for: Products")

	// Retrieve a list of products from the data package
	lp := data.GetProducts()

	// Attempt to convert the list of products to JSON and write it to the response writer
	err := lp.ToJSON(rw)

	// If an error occurs during the conversion, return a 500 Internal Server Error response
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Log a message indicating that a POST request for products was received
	p.l.Println("Handle POST Request for: Products")

	// Create a pointer to a new product
	prod := &data.Product{}

	// Attempt to unmarshal the request body into the product object
	err := prod.FromJSON(r.Body)

	// If an error occurs during the unmarshaling, return a 400 Bad Request response
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	// If the unmarshaling was successful, log the product and add it to the list of products
	if err == nil {
		p.l.Printf("Prod: %#v", prod)
		data.AddProduct(prod)
		createdOn, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", prod.CreatedOn)
		if err != nil {
			log.Fatalf("Error parsing CreatedOn time: %v", err)
		}
		updatedOn, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", prod.UpdatedOn)
		if err != nil {
			log.Fatalf("Error parsing UpdatedOn time: %v", err)
		}
		apple.Exec("INSERT INTO products (id, name, description, price, sku, created_on, updated_on) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			prod.ID, prod.Name, prod.Description, prod.Price, prod.SKU, createdOn, updatedOn)

	}
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Log that a PUT request has been received for products
	p.l.Println("Handle PUT Request for: Products")

	// Create a new product instance
	prod := &data.Product{}

	// Get the ID of the product to be updated from the URL parameters
	idString := ps.ByName("id")

	// Convert the id string to an int
	id, err := strconv.Atoi(idString)
	if err != nil {
		// If the conversion fails, return a Bad Request error
		http.Error(rw, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into the product instance
	err = prod.FromJSON(r.Body)
	if err != nil {
		// If there's an error unmarshaling the request body, return a Bad Request error
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	// Update the product
	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		// If the product is not found, return a Not Found error
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		// If there's an error updating the product, return an Internal Server Error
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

//Example Data for POST: curl -v localhost:9000 -d "{\"id\": 1, \"name\": \"tea\", \"description\": \"a nice cup of tea\"}"
//curl -v localhost:9000 -d "{\"name\": \"tea\", \"description\": \"a nice cup of tea\"}"
