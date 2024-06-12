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
	rows, err := repo.db.Query(`SELECT id, name, description, price FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Get(id int) (*Product, error) {
	row := repo.db.QueryRow(`SELECT id, name, description, price FROM products WHERE id = ?`, id)

	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Create(p Product) (int64, error) {
	result, err := repo.db.Exec(`INSERT INTO products(name, description, price) VALUES (?, ?, ?)`,
		p.Name, p.Description, p.Price)
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
	_, err := repo.db.Exec(`UPDATE products SET name = ?, description = ?, price = ? WHERE id = ?`,
		p.Name, p.Description, p.Price, id)
	return err
}

func (repo *ProductRepository) Delete(id int) error {
	_, err := repo.db.Exec(`DELETE FROM products WHERE id = ?`, id)
	return err
}
