package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Product structure represents a product in the store
type Product struct {
	ID    int
	Name  string
	Size  string
	Price float64
}

// PurchaseRequest structure represents data for the POST request to buy a product
type PurchaseRequest struct {
	ProductID int `json:"product_id"`
	// Add other fields if necessary
}

var products = []Product{
	{1, "Product 1", "Small", 19.99},
	{2, "Product 2", "Medium", 29.99},
	{3, "Product 3", "Large", 39.99},
}

// IndexHandler handles the GET request on the main page of the store
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Online Store</title>
		</head>
		<body>
			<h1>Welcome to the Online Store!</h1>
			<h2>Products:</h2>
			<table border="1">
				<tr>
					<th>ID</th>
					<th>Name</th>
					<th>Size</th>
					<th>Price</th>
					<th>Action</th>
				</tr>
				{{range .}}
					<tr>
						<td>{{.ID}}</td>
						<td>{{.Name}}</td>
						<td>{{.Size}}</td>
						<td>${{.Price}}</td>
						<td><form method="post" action="/buy/{{.ID}}"><input type="submit" value="Buy"></form></td>
					</tr>
				{{end}}
			</table>
			<a href="/add-product">Add Product</a>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, products)
}

// BuyHandler handles the POST request to buy a product
func BuyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var purchaseRequest PurchaseRequest
	if err := decoder.Decode(&purchaseRequest); err != nil {
		http.Error(w, "Invalid JSON message", http.StatusBadRequest)
		return
	}

	// Process data and send a response
	// ...

	fmt.Fprintf(w, "Product with ID %d successfully purchased!", purchaseRequest.ProductID)
}

// AddProductHandler displays the page for adding a new product
func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("addProduct").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Add Product</title>
		</head>
		<body>
			<h1>Add a new product</h1>
			<form method="post" action="/add-product-post">
				<label for="name">Name:</label>
				<input type="text" name="name" required><br>
				<label for="size">Size:</label>
				<input type="text" name="size" required><br>
				<label for="price">Price:</label>
				<input type="number" name="price" step="0.01" required><br>
				<input type="submit" value="Add Product">
			</form>
			<a href="/">Back to Home</a>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// AddProductPostHandler handles the POST request to add a new product
func AddProductPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	size := r.FormValue("size")
	priceStr := r.FormValue("price")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	newProduct := Product{
		ID:    len(products) + 1,
		Name:  name,
		Size:  size,
		Price: price,
	}

	products = append(products, newProduct)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/buy/", BuyHandler)
	http.HandleFunc("/add-product", AddProductHandler)
	http.HandleFunc("/add-product-post", AddProductPostHandler)

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
