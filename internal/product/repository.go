package product

import (
	"fmt"
	"log"
	"strings"
    "github.com/dgomezlikeyoube/ms_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(product *domain.ProductProduct) error
	GetAll(filters Filters, offset int, limit int) ([]domain.Product, error)
	Get(id string) (*domain.Product, error)
	Delete(id string) error
	Update(id string, name *string, sku *string, quantity *int32, price *float32, costPrice *float32, weight *int32, enabled *bool, descripcion *string, category *string) error
	Count(filters Filters) (int, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(product *domain.Product) error {

	if err := repo.db.Create(product).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("El producto fue creado con ID:", product.ID)
	return nil
}

func (repo *repo) GetAll(filters Filters, offset int, limit int) ([]domain.Product, error) {
	var p []domain.Product

	tx := repo.db.Model(&p)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("1").Find(&p)

	if result.Error != nil {
		return nil, result.Error
	}
	return p, nil
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64

	tx := repo.db.Model(domain.Product{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (repo *repo) Get(id string) (*domain.ProductProduct, error) {
	p := domain.ProductProduct{ID: id}
	result := repo.db.First(&p)

	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil

}

func (repo *repo) Delete(id string) error {
	p := domain.Product{ID: id}
	result := repo.db.Delete(&p)

	if result.Error != nil {
		return result.Error
	}
	return nil

}

func (repo *repo) Update(id string, name *string, sku *string, quantity *int32, price *float32, costPrice *float32, weight *int32, enabled *bool, descripcion *string, category *string) error {
	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if sku != nil {
		values["quantity"] = *quantity
	}

	if price != nil {
		values["price"] = *price
	}

	if costPrice != nil {
		values["costPrice"] = *costPrice
	}

	if weight != nil {
		values["weight"] = *weight
	}

	if enabled != nil {
		values["enabled"] = *enabled
	}

	if descripcion != nil {
		values["descripcion"] = *descripcion
	}

	if category != nil {
		values["category"] = *category
	}

	println(id)
	if err := repo.db.Model(&domain.Product{}).Where("id=?", id).Updates(values).Error; err != nil {
		return err
	}

	return nil

}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Sku != "" {
		filters.Sku = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Sku))
		tx = tx.Where("lower(sku) like ?", filters.Sku)
	}

	return tx
}
