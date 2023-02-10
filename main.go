package main

import (
	"database/sql"
	"log"
	"main/work/data"
	"main/work/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func main() {
	db, err := initDatabase()
	if err != nil {
		log.Println("Error initializing the database")
	}
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
	ph := handlers.NewProducts(l, db)
	router := httprouter.New()
	router.GET("/products", ph.GetProducts)
	router.POST("/products", ph.AddProduct)
	router.PUT("/products/:id", ph.UpdateProduct)
	log.Fatal(http.ListenAndServe(":8080", router))
	//sm := http.NewServeMux()
	//sm.Handle("/", ph)

	/*
		s := &http.Server{
			Addr: ":9000",
			//Handler:      sm,
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
	*/

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, gracefully shuttingdown.\n Signal Type: ", sig)
	//tc, cf := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cf()
	//s.Shutdown(tc)
}

func initDatabase() (*sql.DB, error) {
	// A flag to control adding a few extra logs for debugging purposes
	test := false
	// Open a connection to the "product_db" database using the PostgreSQL driver, with the specified user credentials
	var db, err = sql.Open("postgres", "user=postgres password=1 dbname=product_db sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	// Defer the closing of the database connection until the function returns
	defer db.Close()

	// Insert product data
	for _, product := range data.GetProducts() {

		// Check if the product already exists in the database by querying the "products" table for a matching id or sku
		var id int
		err := db.QueryRow("SELECT id FROM products WHERE id = $1 OR sku = $2", product.ID, product.SKU).Scan(&id)
		if err != nil {
			// If there's no match, parse the "created_on" and "updated_on" timestamps into Go `time.Time` values
			createdOn, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", product.CreatedOn)
			if err != nil {
				log.Fatalf("Error parsing CreatedOn time: %v", err)
			}
			updatedOn, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", product.UpdatedOn)
			if err != nil {
				log.Fatalf("Error parsing UpdatedOn time: %v", err)
			}
			// Execute an INSERT statement to insert the product data into the "products" table
			_, err = db.Exec("INSERT INTO products (id, name, description, price, sku, created_on, updated_on) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				product.ID, product.Name, product.Description, product.Price, product.SKU, createdOn, updatedOn)
			if err != nil {
				log.Fatalf("Error inserting product data: %v", err)
			}
		} else {
			if test {
				// If the "test" flag is set to true, log a message indicating that the product already exists in the database
				log.Printf("Attempted to added product(s). Product(s) already exist in the DB: %#v", product)
			}
		}

	}
	return db, nil
}
