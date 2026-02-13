package services

import (
	"errors"
	"order-management-system/internal/domain"
	"testing"
)

// Mock Repositories
type mockUserRepository struct {
	users map[uint]*domain.User
}

func (m *mockUserRepository) GetByID(id uint) (*domain.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepository) Create(user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	for _, u := range m.users {
		users = append(users, *u)
	}
	return users, nil
}

type mockProductRepository struct {
	products map[uint]*domain.Product
}

func (m *mockProductRepository) GetByID(id uint) (*domain.Product, error) {
	if product, ok := m.products[id]; ok {
		return product, nil
	}
	return nil, errors.New("product not found")
}

func (m *mockProductRepository) UpdateStock(id uint, quantity int) error {
	if product, ok := m.products[id]; ok {
		product.Stock = quantity
		return nil
	}
	return errors.New("product not found")
}

func (m *mockProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	for _, p := range m.products {
		products = append(products, *p)
	}
	return products, nil
}

func (m *mockProductRepository) Create(product *domain.Product) error {
	m.products[product.ID] = product
	return nil
}

type mockOrderRepository struct {
	orders map[uint]*domain.Order
	nextID uint
}

func (m *mockOrderRepository) Create(order *domain.Order) error {
	m.nextID++
	order.ID = m.nextID
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) GetByID(id uint) (*domain.Order, error) {
	if order, ok := m.orders[id]; ok {
		return order, nil
	}
	return nil, errors.New("order not found")
}

func (m *mockOrderRepository) GetAll() ([]domain.Order, error) {
	var orders []domain.Order
	for _, o := range m.orders {
		orders = append(orders, *o)
	}
	return orders, nil
}

func (m *mockOrderRepository) GetByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	for _, o := range m.orders {
		if o.UserID == userID {
			orders = append(orders, *o)
		}
	}
	return orders, nil
}

func (m *mockOrderRepository) Update(order *domain.Order) error {
	if _, ok := m.orders[order.ID]; ok {
		m.orders[order.ID] = order
		return nil
	}
	return errors.New("order not found")
}

// Test Functions
func setupService() (*OrderService, *mockUserRepository, *mockProductRepository, *mockOrderRepository) {
	userRepo := &mockUserRepository{users: make(map[uint]*domain.User)}
	productRepo := &mockProductRepository{products: make(map[uint]*domain.Product)}
	orderRepo := &mockOrderRepository{orders: make(map[uint]*domain.Order), nextID: 0}

	// Setup test data
	userRepo.users[1] = &domain.User{ID: 1, Name: "Test User", Email: "test@test.com"}
	productRepo.products[1] = &domain.Product{ID: 1, Name: "Product 1", Price: 100.0, Stock: 10}
	productRepo.products[2] = &domain.Product{ID: 2, Name: "Product 2", Price: 50.0, Stock: 5}

	service := NewOrderService(orderRepo, productRepo, userRepo)
	return service, userRepo, productRepo, orderRepo
}

func TestCreateOrder_Success(t *testing.T) {
	service, _, _, _ := setupService()

	req := domain.CreateOrderRequest{
		UserID: 1,
		Items: []domain.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	}

	order, err := service.CreateOrder(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if order.Status != domain.StatusPending {
		t.Errorf("Expected status PENDING, got %s", order.Status)
	}

	expectedTotal := 2000.0 // 2 * 100
	if order.Total != expectedTotal {
		t.Errorf("Expected total %f, got %f", expectedTotal, order.Total)
	}
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	service, _, _, _ := setupService()

	req := domain.CreateOrderRequest{
		UserID: 1,
		Items: []domain.OrderItemRequest{
			{ProductID: 1, Quantity: 20}, // Stock is only 10
		},
	}

	_, err := service.CreateOrder(req)
	if err != ErrInsufficientStock {
		t.Errorf("Expected ErrInsufficientStock, got %v", err)
	}
}

func TestCreateOrder_UserNotFound(t *testing.T) {
	service, _, _, _ := setupService()

	req := domain.CreateOrderRequest{
		UserID: 999, // Non-existent user
		Items: []domain.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	}

	_, err := service.CreateOrder(req)
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestConfirmOrder_Success(t *testing.T) {
	service, _, productRepo, _ := setupService()

	// Create order first
	req := domain.CreateOrderRequest{
		UserID: 1,
		Items: []domain.OrderItemRequest{
			{ProductID: 1, Quantity: 3},
		},
	}
	order, _ := service.CreateOrder(req)

	initialStock := productRepo.products[1].Stock

	// Confirm order
	confirmedOrder, err := service.ConfirmOrder(order.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if confirmedOrder.Status != domain.StatusConfirmed {
		t.Errorf("Expected status CONFIRMED, got %s", confirmedOrder.Status)
	}

	// Check stock reduction
	expectedStock := initialStock - 1
	if productRepo.products[1].Stock != expectedStock {
		t.Errorf("Expected stock %d, got %d", expectedStock, productRepo.products[1].Stock)
	}
}

func TestConfirmOrder_InvalidStatus(t *testing.T) {
	service, _, _, _ := setupService()

	// Create and confirm order
	req := domain.CreateOrderRequest{
		UserID: 1,
		Items: []domain.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	}
	order, _ := service.CreateOrder(req)
	service.ConfirmOrder(order.ID)

	// Try to confirm again
	_, err := service.ConfirmOrder(order.ID)
	if err != ErrInvalidStatus {
		t.Errorf("Expected ErrInvalidStatus, got %v", err)
	}
}

func TestShipOrder_Success(t *testing.T) {
	service, _, _, _ := setupService()

	req := domain.CreateOrderRequest{
		UserID: 1,
		Items: []domain.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	}
	order, _ := service.CreateOrder(req)
	service.ConfirmOrder(order.ID)

	shippedOrder, err := service.ShipOrder(order.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if shippedOrder.Status != domain.StatusShipped {
		t.Errorf("Expected status SHIPPED, got %s", shippedOrder.Status)
	}
}

func TestCancelOrder_ReturnStock(t *testing.T) {
	service, _, productRepo, _ := setupService()

	req := domain.CreateOrderRequest{
		UserID: 1,
		Items: []domain.OrderItemRequest{
			{ProductID: 1, Quantity: 3},
		},
	}
	order, _ := service.CreateOrder(req)
	service.ConfirmOrder(order.ID)

	stockAfterConfirm := productRepo.products[1].Stock

	// Cancel order
	cancelledOrder, err := service.CancelOrder(order.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cancelledOrder.Status != domain.StatusCancelled {
		t.Errorf("Expected status CANCELLED, got %s", cancelledOrder.Status)
	}

	// Check stock returned
	expectedStock := stockAfterConfirm + 3
	if productRepo.products[1].Stock != expectedStock {
		t.Errorf("Expected stock %d, got %d", expectedStock, productRepo.products[1].Stock)
	}
}

func TestCancelOrder_CannotCancelShipped(t *testing.T) {
	service, _, _, _ := setupService()

	req := domain.CreateOrderRequest{
		UserID: 1,
		Items: []domain.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	}
	order, _ := service.CreateOrder(req)
	service.ConfirmOrder(order.ID)
	service.ShipOrder(order.ID)

	_, err := service.CancelOrder(order.ID)
	if err != ErrCannotCancelShipped {
		t.Errorf("Expected ErrCannotCancelShipped, got %v", err)
	}
}
