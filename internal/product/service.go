package product

import (
	"log"
	 "github.com/dgomezlikeyoube/ms_domain/domain"
)

type (
	Filters struct {
		Sku string
	}

	Service interface {
		Create(name string, sku string, quantity int32, price float32, costprice float32, wight int32, enabled bool, descripcion string, category string) (*domain.Product, error)
		Get(id string) (*domain.Product, error)
		GetAll(filters Filters, offsert int, limit int) ([]domain.Product, error)
		Delete(id string) error
		Update(id string, name *string, sku *string, quantity *int32, price *float32, costPrice *float32, weight *int32, enabled *bool, descripcion *string, category *string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(name string, sku string, quantity int32, price float32, costprice float32, wight int32, enabled bool, descripcion string, category string) (*domain.Product, error) {
	s.log.Println("Create product service")
	product := domain.Product{
		Name:        name,
		Sku:         sku,
		Quantity:    quantity,
		Price:       price,
		CostPrice:   costprice,
		Weight:      wight,
		Enabled:     enabled,
		Descripcion: descripcion,
		Category:    category,
	}
	if err := s.repo.Create(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s service) GetAll(filters Filters, offset int, limit int) ([]domain.Product, error) {
	product, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s service) Get(id string) (*domain.Product, error) {
	product, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Update(id string, name *string, sku *string, quantity *int32, price *float32, costPrice *float32, weight *int32, enabled *bool, descripcion *string, category *string) error {
	println(*name)
	return s.repo.Update(id, name, sku, quantity, price, costPrice, weight, enabled, descripcion, category)

}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)

}
