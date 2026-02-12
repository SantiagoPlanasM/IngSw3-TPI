package main

import (
	"log"
	"order-management-system/internal/config"
	"order-management-system/internal/handlers"
	"order-management-system/internal/repositories"
	"order-management-system/internal/services"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Seed database with initial data
	if err := config.SeedDatabase(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Initialize services
	orderService := services.NewOrderService(orderRepo, productRepo, userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)
	productHandler := handlers.NewProductHandler(productRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup Gin router
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:80", "http://localhost:5173", "http://localhost:3000", "http://frontend:80"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.POST("", userHandler.Create)
		}

		// Product routes
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
			products.POST("", productHandler.Create)
		}

		// Order routes
		orders := api.Group("/orders")
		{
			orders.GET("", orderHandler.GetAll)
			orders.GET("/:id", orderHandler.GetByID)
			orders.GET("/user/:userId", orderHandler.GetByUserID)
			orders.POST("", orderHandler.Create)
			orders.PATCH("/:id/confirm", orderHandler.Confirm)
			orders.PATCH("/:id/ship", orderHandler.Ship)
			orders.PATCH("/:id/cancel", orderHandler.Cancel)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
