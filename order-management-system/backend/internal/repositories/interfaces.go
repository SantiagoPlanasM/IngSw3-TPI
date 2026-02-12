package repositories

import "order-management-system/internal/domain"

type UserRepository interface {
	GetByID(id uint) (*domain.User, error)
	Create(user *domain.User) error
	GetAll() ([]domain.User, error)
}

type ProductRepository interface {
	GetByID(id uint) (*domain.Product, error)
	UpdateStock(id uint, quantity int) error
	GetAll() ([]domain.Product, error)
	Create(product *domain.Product) error
}

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id uint) (*domain.Order, error)
	GetAll() ([]domain.Order, error)
	GetByUserID(userID uint) ([]domain.Order, error)
	Update(order *domain.Order) error
}
