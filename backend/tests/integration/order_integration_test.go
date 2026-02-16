package integration

import (
	"order-management-system/internal/config"
	"order-management-system/internal/domain"
	"order-management-system/internal/repositories"
	"order-management-system/internal/services"
	"os"
	"testing"
)

// This integration test requires a real MySQL database
// Set environment variables before running:
// DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME

func TestOrderIntegration(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Clean up after test
	defer func() {
		db.Exec("DELETE FROM order_items")
		db.Exec("DELETE FROM orders")
		db.Exec("DELETE FROM products")
		db.Exec("DELETE FROM users")
	}()

	// Setup repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Setup service
	orderService := services.NewOrderService(orderRepo, productRepo, userRepo)

	// Create test user
	user := &domain.User{
		Name:  "Integration Test User",
		Email: "integration@test.com",
	}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create test product
	product := &domain.Product{
		Name:  "Test Product",
		Price: 100.0,
		Stock: 10,
	}
	if err := productRepo.Create(product); err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Test: Create Order
	t.Run("Create Order", func(t *testing.T) {
		req := domain.CreateOrderRequest{
			UserID: user.ID,
			Items: []domain.OrderItemRequest{
				{ProductID: product.ID, Quantity: 2},
			},
		}

		order, err := orderService.CreateOrder(req)
		if err != nil {
			t.Fatalf("Failed to create order: %v", err)
		}

		if order.Status != domain.StatusPending {
			t.Errorf("Expected status PENDING, got %s", order.Status)
		}

		expectedTotal := 200.0
		if order.Total != expectedTotal {
			t.Errorf("Expected total %f, got %f", expectedTotal, order.Total)
		}

		// Test: Confirm Order
		t.Run("Confirm Order", func(t *testing.T) {
			confirmedOrder, err := orderService.ConfirmOrder(order.ID)
			if err != nil {
				t.Fatalf("Failed to confirm order: %v", err)
			}

			if confirmedOrder.Status != domain.StatusConfirmed {
				t.Errorf("Expected status CONFIRMED, got %s", confirmedOrder.Status)
			}

			// Verify stock reduction
			updatedProduct, err := productRepo.GetByID(product.ID)
			if err != nil {
				t.Fatalf("Failed to get product: %v", err)
			}

			expectedStock := 80 // 10 - 2
			if updatedProduct.Stock != expectedStock {
				t.Errorf("Expected stock %d, got %d", expectedStock, updatedProduct.Stock)
			}

			// Test: Ship Order
			t.Run("Ship Order", func(t *testing.T) {
				shippedOrder, err := orderService.ShipOrder(order.ID)
				if err != nil {
					t.Fatalf("Failed to ship order: %v", err)
				}

				if shippedOrder.Status != domain.StatusShipped {
					t.Errorf("Expected status SHIPPED, got %s", shippedOrder.Status)
				}
			})
		})
	})

	// Test: Cancel Order with Stock Return
	t.Run("Cancel Order Returns Stock", func(t *testing.T) {
		// Create new order
		req := domain.CreateOrderRequest{
			UserID: user.ID,
			Items: []domain.OrderItemRequest{
				{ProductID: product.ID, Quantity: 3},
			},
		}

		order, err := orderService.CreateOrder(req)
		if err != nil {
			t.Fatalf("Failed to create order: %v", err)
		}

		// Confirm order (reduces stock)
		_, err = orderService.ConfirmOrder(order.ID)
		if err != nil {
			t.Fatalf("Failed to confirm order: %v", err)
		}

		stockBeforeCancel, _ := productRepo.GetByID(product.ID)

		// Cancel order
		cancelledOrder, err := orderService.CancelOrder(order.ID)
		if err != nil {
			t.Fatalf("Failed to cancel order: %v", err)
		}

		if cancelledOrder.Status != domain.StatusCancelled {
			t.Errorf("Expected status CANCELLED, got %s", cancelledOrder.Status)
		}

		// Verify stock returned
		stockAfterCancel, _ := productRepo.GetByID(product.ID)
		expectedStock := stockBeforeCancel.Stock + 3

		if stockAfterCancel.Stock != expectedStock {
			t.Errorf("Expected stock %d, got %d", expectedStock, stockAfterCancel.Stock)
		}
	})
}
