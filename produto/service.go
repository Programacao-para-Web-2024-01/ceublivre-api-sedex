package produto

type ProductService struct {
	repository *ProductRepository
}

func NewProductService(repository *ProductRepository) *ProductService {
	return &ProductService{repository: repository}
}

func (svc *ProductService) List() ([]Product, error) {
	return svc.repository.List()
}

func (svc *ProductService) Get(id int) (*Product, error) {
	return svc.repository.Get(id)
}

func (svc *ProductService) Create(p Product) (*Product, error) {
	newId, err := svc.repository.Create(p)
	if err != nil {
		return nil, err
	}

	p.ID = newId
	return &p, nil
}

func (svc *ProductService) Update(p Product) error {
	id := int(p.ID)

	_, err := svc.Get(id)
	if err != nil {
		return err
	}

	return svc.repository.Update(id, p)
}

func (svc *ProductService) Delete(id int) error {
	_, err := svc.Get(id)
	if err != nil {
		return err
	}

	return svc.repository.Delete(id)
}

func (svc *ProductService) IsAvailable(id int, quantity int) (bool, error) {
	product, err := svc.Get(id)
	if err != nil {
		return false, err
	}
	return product.QuantityStock >= quantity, nil
}
