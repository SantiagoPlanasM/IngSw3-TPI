package services

import (
	"errors"
	"order-management-system/internal/domain"
	"order-management-system/internal/repositories"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidStatus      = errors.New("invalid order status transition")
	ErrCannotCancelShipped = errors.New("cannot cancel shipped order")
)

type OrderService struct {
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
	userRepo    repositories.UserRepository
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	userRepo repositories.UserRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// CreateOrder valida stock, existencia de usuario, calcula total y crea pedido con estado PENDING
func (s *OrderService) CreateOrder(req domain.CreateOrderRequest) (*domain.Order, error) {
	// Validar existencia del usuario
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	var total float64
	var orderItems []domain.OrderItem

	// Validar stock y calcular total
	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, ErrProductNotFound
		}

		if product.Stock < item.Quantity {
			return nil, ErrInsufficientStock
		}

		orderItem := domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		orderItems = append(orderItems, orderItem)
		total += product.Price * float64(item.Quantity)
	}

	order := &domain.Order{
		UserID: user.ID,
		Total:  total,
		Status: domain.StatusPending,
		Items:  orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(order.ID)
}

// ConfirmOrder reduce el stock real y cambia el estado a CONFIRMED
func (s *OrderService) ConfirmOrder(orderID uint) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.Status != domain.StatusPending {
		return nil, ErrInvalidStatus
	}

	// Reducir stock de cada producto
	for _, item := range order.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, ErrProductNotFound
		}

		newStock := product.Stock - item.Quantity
		if newStock < 0 {
			return nil, ErrInsufficientStock
		}

		if err := s.productRepo.UpdateStock(item.ProductID, newStock); err != nil {
			return nil, err
		}
	}

	order.Status = domain.StatusConfirmed
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(order.ID)
}

// ShipOrder cambia el estado a SHIPPED solo si está CONFIRMED
func (s *OrderService) ShipOrder(orderID uint) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.Status != domain.StatusConfirmed {
		return nil, ErrInvalidStatus
	}

	order.Status = domain.StatusShipped
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(order.ID)
}

// CancelOrder devuelve el stock si no fue enviado y cambia estado a CANCELLED
func (s *OrderService) CancelOrder(orderID uint) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.Status == domain.StatusShipped {
		return nil, ErrCannotCancelShipped
	}

	// Si el pedido estaba confirmado, devolver stock
	if order.Status == domain.StatusConfirmed {
		for _, item := range order.Items {
			product, err := s.productRepo.GetByID(item.ProductID)
			if err != nil {
				return nil, ErrProductNotFound
			}

			newStock := product.Stock + item.Quantity
			if err := s.productRepo.UpdateStock(item.ProductID, newStock); err != nil {
				return nil, err
			}
		}
	}

	order.Status = domain.StatusCancelled
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(order.ID)
}

func (s *OrderService) GetOrder(orderID uint) (*domain.Order, error) {
	return s.orderRepo.GetByID(orderID)
}

func (s *OrderService) GetAllOrders() ([]domain.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *OrderService) GetOrdersByUser(userID uint) ([]domain.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}
