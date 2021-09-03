package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/chrismeh/lefty/internal/inmem"
	"github.com/chrismeh/lefty/pkg/retailer"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	infoLog      *log.Logger
	errorLog     *log.Logger
	productStore *inmem.ProductStore
}

func main() {
	addr := flag.String("port", ":5000", "HTTP address to listen on")
	flag.Parse()

	app := application{
		infoLog:      log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog:     log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		productStore: inmem.NewProductStore(),
	}

	router := http.NewServeMux()
	router.HandleFunc("/api/products", app.handleGetProducts)

	s := &http.Server{
		Addr:         *addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	go app.updateRetailers()

	app.infoLog.Printf("starting application at %s", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func (a application) json(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(v)
}

func (a application) jsonError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": error})
}

func (a application) updateRetailers() {
	f, err := os.Open("products.json")
	if err == nil {
		a.infoLog.Println("Skipped retailer update: products.json found")
		_ = a.productStore.Load(f)
		f.Close()
		return
	}
	if !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	defer f.Close()

	start := time.Now()
	a.infoLog.Println("Starting retailer update ...")

	c := http.Client{Timeout: 5 * time.Second}
	err = retailer.UpdateRetailers(a.productStore, retailer.NewThomann(&c), retailer.NewMusikProduktiv(&c))
	if err != nil {
		a.errorLog.Println(err)
	}

	duration := time.Since(start)
	a.infoLog.Printf("Finished retailer update after %d ms", duration.Milliseconds())

	f, err = os.Create("products.json")
	if err != nil {
		panic(err)
	}

	_ = a.productStore.Dump(f)
	f.Close()
}
