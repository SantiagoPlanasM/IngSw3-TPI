package repositories

import (
	"order-management-system/internal/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id uint) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.Preload("User").Preload("Items.Product").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetAll() ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Preload("User").Preload("Items.Product").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Where("user_id = ?", userID).Preload("User").Preload("Items.Product").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}
