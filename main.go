package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB


type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    float64 `json:"price"`
}


func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(mysql:3306)/products_db")
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
}


func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO products (name, category, price) VALUES (?, ?, ?)"
	res, err := db.Exec(query, product.Name, product.Category, product.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	product.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	row := db.QueryRow("SELECT id, name, category, price FROM products WHERE id = ?", id)

	var product Product
	if err := row.Scan(&product.ID, &product.Name, &product.Category, &product.Price); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}


func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET name = ?, category = ?, price = ? WHERE id = ?"
	_, err := db.Exec(query, product.Name, product.Category, product.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func listProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, category, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}


func listProductsByCategory(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["category"]
	rows, err := db.Query("SELECT id, name, category, price FROM products WHERE category = ?", category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}


func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", getProduct).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", updateProduct).Methods("PUT")
	r.HandleFunc("/products", listProducts).Methods("GET")
	r.HandleFunc("/products/category/{category}", listProductsByCategory).Methods("GET")

	fmt.Println("Server is running on port 0.0.0.0:3004")
	log.Fatal(http.ListenAndServe("0.0.0.0:3004", r))
}
