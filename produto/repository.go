package produto

import (
	"database/sql"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(database *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: database,
	}
}

func (repo *ProductRepository) List() ([]Product, error) {
	rows, err := repo.db.Query(`SELECT product_id, name, description, price, quantity_stock FROM ProductCatalog`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.QuantityStock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Get(id int) (*Product, error) {
	row := repo.db.QueryRow(`SELECT product_id, name, description, price, quantity_stock FROM ProductCatalog WHERE product_id = ?`, id)

	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.QuantityStock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Create(p Product) (int64, error) {
	result, err := repo.db.Exec(`INSERT INTO ProductCatalog(name, description, price, quantity_stock) VALUES (?, ?, ?, ?)`,
		p.Name, p.Description, p.Price, p.QuantityStock)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *ProductRepository) Update(id int, p Product) error {
	_, err := repo.db.Exec(`UPDATE ProductCatalog SET name = ?, description = ?, price = ?, quantity_stock = ? WHERE product_id = ?`,
		p.Name, p.Description, p.Price, p.QuantityStock, id)
	return err
}

func (repo *ProductRepository) Delete(id int) error {
	_, err := repo.db.Exec(`DELETE FROM ProductCatalog WHERE product_id = ?`, id)
	return err
}
