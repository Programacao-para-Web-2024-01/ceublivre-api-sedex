package carrinho

type Cart struct {
	ID        int64      `json:"id"`
	CreatedAt string     `json:"created_at"`
	Items     []CartItem `json:"items,omitempty"`
	Active    bool       `json:"active"`
}

type CartItem struct {
	ID        int64 `json:"id"`
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
