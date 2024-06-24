package carrinho

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type CartController struct {
	svc *CartService
}

func NewCartController(svc *CartService) *CartController {
	return &CartController{svc: svc}
}

func (ctrl *CartController) AddItem(w http.ResponseWriter, r *http.Request) {
	var item CartItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Recebida requisição POST para adicionar item:", item)

	newItem, err := ctrl.svc.AddItem(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctrl *CartController) RemoveItem(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctrl.svc.RemoveItem(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *CartController) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item CartItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctrl.svc.UpdateItem(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *CartController) CalculateTotal(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	total, err := ctrl.svc.CalculateTotal(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]float64{"total": total})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *CartController) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var items []CartItem
	err := json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	available, err := ctrl.svc.CheckAvailability(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]bool{"available": available})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *CartController) CreateCart(w http.ResponseWriter, r *http.Request) {
	var cart Cart
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCart, err := ctrl.svc.Create(cart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(newCart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctrl *CartController) GetActiveCart(w http.ResponseWriter, r *http.Request) {
	cart, err := ctrl.svc.GetActiveCart()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(cart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterCartRoutes(router *mux.Router, svc *CartService) {
	ctrl := NewCartController(svc)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	cartRouter := router.PathPrefix("/cart").Subrouter()
	cartRouter.Use(cors)

	cartRouter.HandleFunc("/items", ctrl.AddItem).Methods("POST")
	cartRouter.HandleFunc("/items/{id}", ctrl.RemoveItem).Methods("DELETE")
	cartRouter.HandleFunc("/items/{id}", ctrl.UpdateItem).Methods("PUT")
	cartRouter.HandleFunc("/total/{id}", ctrl.CalculateTotal).Methods("GET")
	cartRouter.HandleFunc("/availability", ctrl.CheckAvailability).Methods("POST")
	cartRouter.HandleFunc("", ctrl.CreateCart).Methods("POST")
	cartRouter.HandleFunc("/active", ctrl.GetActiveCart).Methods("GET")
}
