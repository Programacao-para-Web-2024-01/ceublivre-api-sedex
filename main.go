package main

import (
	"api-carrinho/carrinho"
	"api-carrinho/produto"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	if err := createServer(); err != nil {
		log.Panic(err)
	}
}

func connectDB() *sql.DB {
	config := mysql.NewConfig()
	config.User = "root"
	config.Passwd = "uniceub"
	config.DBName = "shop"
	conn, err := mysql.NewConnector(config)
	if err != nil {
		panic(err)
	}
	return sql.OpenDB(conn)
}

func createServer() error {
	db := connectDB()

	productRepository := produto.NewProductRepository(db)
	cartRepository := carrinho.NewCartRepository(db)

	productSvc := produto.NewProductService(productRepository)
	cartSvc := carrinho.NewCartService(cartRepository, productSvc)

	productController := produto.NewProductController(productSvc)
	cartController := carrinho.NewCartController(cartSvc)

	router := mux.NewRouter()

	router.HandleFunc("/products", productController.List).Methods("GET")
	router.HandleFunc("/products/{id}", productController.Get).Methods("GET")
	router.HandleFunc("/products", productController.Create).Methods("POST")
	router.HandleFunc("/products/{id}", productController.Update).Methods("PUT")
	router.HandleFunc("/products/{id}", productController.Delete).Methods("DELETE")

	router.HandleFunc("/cart", cartController.GetActiveCart).Methods("GET")
	router.HandleFunc("/cart", cartController.CreateCart).Methods("POST")
	router.HandleFunc("/cart/items", cartController.AddItem).Methods("POST")
	router.HandleFunc("/cart/items/{id}", cartController.UpdateItem).Methods("PUT")
	router.HandleFunc("/cart/items/{id}", cartController.RemoveItem).Methods("DELETE")
	router.HandleFunc("/cart/total", cartController.CalculateTotal).Methods("GET")
	router.HandleFunc("/cart/availability", cartController.CheckAvailability).Methods("POST")

	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Println("Servidor rodando em http://localhost:8080")
	return http.ListenAndServe(":8080", corsOptions(router))
}
