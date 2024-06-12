package carrinho

import (
	"database/sql"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(database *sql.DB) *CartRepository {
	return &CartRepository{
		db: database,
	}
}

func (repo *CartRepository) List() ([]Cart, error) {
	rows, err := repo.db.Query(`SELECT id, created_at FROM carts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []Cart
	for rows.Next() {
		var c Cart
		err = rows.Scan(&c.ID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}

		items, err := repo.ListItems(c.ID)
		if err != nil {
			return nil, err
		}
		c.Items = items

		carts = append(carts, c)
	}

	return carts, nil
}

func (repo *CartRepository) Get(id int64) (*Cart, error) {
	row := repo.db.QueryRow(`SELECT id, created_at FROM carts WHERE id = ?`, id)

	var c Cart
	err := row.Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	items, err := repo.ListItems(c.ID)
	if err != nil {
		return nil, err
	}
	c.Items = items

	return &c, nil
}

func (repo *CartRepository) Create(c Cart) (int64, error) {
	result, err := repo.db.Exec(`INSERT INTO carts(created_at) VALUES (?)`, c.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *CartRepository) Update(id int64, c Cart) error {
	_, err := repo.db.Exec(`UPDATE carts SET created_at = ? WHERE id = ?`, c.CreatedAt, id)
	return err
}

func (repo *CartRepository) Delete(id int64) error {
	_, err := repo.db.Exec(`DELETE FROM carts WHERE id = ?`, id)
	return err
}

func (repo *CartRepository) ListItems(cartID int64) ([]CartItem, error) {
	rows, err := repo.db.Query(`SELECT id, cart_id, product_id, quantity FROM cart_items WHERE cart_id = ?`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		err = rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (repo *CartRepository) AddItem(item CartItem) (int64, error) {
	result, err := repo.db.Exec(`INSERT INTO cart_items(cart_id, product_id, quantity) VALUES (?, ?, ?)`, item.CartID, item.ProductID, item.Quantity)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *CartRepository) UpdateItem(item CartItem) error {
	_, err := repo.db.Exec(`UPDATE cart_items SET quantity = ? WHERE id = ?`, item.Quantity, item.ID)
	return err
}

func (repo *CartRepository) RemoveItem(id int64) error {
	_, err := repo.db.Exec(`DELETE FROM cart_items WHERE id = ?`, id)
	return err
}
