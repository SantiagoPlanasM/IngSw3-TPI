package repositories

import (
	"order-management-system/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&domain.Product{}).Where("id = ?", id).Update("stock", quantity).Error
}

func (r *productRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}
