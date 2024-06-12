package main

import (
	"api-carrinho/carrinho"
	"api-carrinho/produto"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
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
	config.DBName = "web"
	conn, err := mysql.NewConnector(config)
	if err != nil {
		panic(err)
	}
	return sql.OpenDB(conn)
}

func createServer() error {
	db := connectDB()

	ProductRepository := produto.NewProductRepository(db)
	CartRepository := carrinho.NewCartRepository(db)
	productSvc := produto.NewProductService(ProductRepository)
	cartSvc := carrinho.NewCartService(CartRepository, productSvc)
	productController := produto.NewProductController(productSvc)
	cartController := carrinho.NewCartController(cartSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products/", productController.List)
	mux.HandleFunc("GET /products/{id}", productController.Get)
	mux.HandleFunc("POST /products/", productController.Create)
	mux.HandleFunc("PUT /products/{id}", productController.Update)
	mux.HandleFunc("DELETE /products/{id}", productController.Delete)

	mux.HandleFunc("GET /cart", cartController.GetActiveCart)
	mux.HandleFunc("POST /cart", cartController.CreateCart)
	mux.HandleFunc("POST /cart/item", cartController.AddItem)
	mux.HandleFunc("PUT /cart/item/{id}", cartController.UpdateItem)
	mux.HandleFunc("DELETE /cart/item/{id}", cartController.RemoveItem)
	mux.HandleFunc("GET /cart/total", cartController.CalculateTotal)
	mux.HandleFunc("POST /cart/availability", cartController.CheckAvailability)

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

	return http.ListenAndServe("localhost:8080", mux)
}

func appendMiddlewares(
	handler func(w http.ResponseWriter, req *http.Request),
	mw ...func(w http.ResponseWriter, req *http.Request) error,
) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, middleware := range mw {
			err := middleware(w, req)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}
		}

		handler(w, req)
	}
}

func authentication(w http.ResponseWriter, req *http.Request) error {
	authorization := req.Header.Get("Authorization")
	_, err := validateToken(authorization)
	if err != nil {
		w.WriteHeader(401)
		return err
	}

	return nil
}

var key = []byte("TOKEN_SECRETO")
var jwtManager = jwt.New(jwt.SigningMethodHS256)

func createToken() (string, error) {
	return jwtManager.SignedString(key)
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.NewParser().Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
}
