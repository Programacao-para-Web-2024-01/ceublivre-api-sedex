package carrinho

import (
	"api-carrinho/produto"
	"errors"
)

type CartService struct {
	repository *CartRepository
	productSvc *produto.ProductService
}

func NewCartService(repository *CartRepository, productSvc *produto.ProductService) *CartService {
	return &CartService{repository: repository, productSvc: productSvc}
}

func (svc *CartService) List() ([]Cart, error) {
	return svc.repository.List()
}

func (svc *CartService) Get(id int64) (*Cart, error) {
	return svc.repository.Get(id)
}

func (svc *CartService) Create(c Cart) (*Cart, error) {
	newId, err := svc.repository.Create(c)
	if err != nil {
		return nil, err
	}

	c.ID = newId
	return &c, nil
}

func (svc *CartService) Update(c Cart) error {
	id := c.ID

	_, err := svc.Get(id)
	if err != nil {
		return err
	}

	return svc.repository.Update(id, c)
}

func (svc *CartService) Delete(id int64) error {
	_, err := svc.Get(id)
	if err != nil {
		return err
	}

	return svc.repository.Delete(id)
}

func (svc *CartService) AddItem(item CartItem) (*CartItem, error) {
	available, err := svc.productSvc.IsAvailable(int(item.ProductID), item.Quantity)
	if err != nil {
		return nil, err
	}

	if !available {
		return nil, errors.New("produto não disponível em estoque")
	}

	newId, err := svc.repository.AddItem(item)
	if err != nil {
		return nil, err
	}

	item.ID = newId
	return &item, nil
}

func (svc *CartService) UpdateItem(item CartItem) error {
	available, err := svc.productSvc.IsAvailable(int(item.ProductID), item.Quantity)
	if err != nil {
		return err
	}

	if !available {
		return errors.New("produto não disponível em estoque")
	}

	return svc.repository.UpdateItem(item)
}

func (svc *CartService) RemoveItem(id int64) error {
	return svc.repository.RemoveItem(id)
}

func (svc *CartService) CalculateTotal(cartID int64) (float64, error) {
	cart, err := svc.Get(cartID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range cart.Items {
		product, err := svc.productSvc.Get(int(item.ProductID))
		if err != nil {
			return 0, err
		}
		total += product.Price * float64(item.Quantity)
	}

	return total, nil
}

func (svc *CartService) CheckAvailability(items []CartItem) (bool, error) {
	for _, item := range items {
		product, err := svc.productSvc.Get(int(item.ProductID))
		if err != nil {
			return false, err
		}
		if product.QuantityStock < item.Quantity {
			return false, nil
		}
	}
	return true, nil
}

func (svc *CartService) GetActiveCart() (*Cart, error) {
	carts, err := svc.List()
	if err != nil {
		return nil, err
	}

	for _, cart := range carts {
		if cart.Active {
			return &cart, nil
		}
	}

	newCart, err := svc.Create(Cart{})
	if err != nil {
		return nil, err
	}

	return newCart, nil
}
