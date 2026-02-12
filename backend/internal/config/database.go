package config

import (
	"fmt"
	"log"
	"order-management-system/internal/domain"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Order{},
		&domain.OrderItem{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}

func SeedDatabase(db *gorm.DB) error {
	// Check if data already exists
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already seeded")
		return nil
	}

	// Seed users
	users := []domain.User{
		{Name: "Juan Pérez", Email: "juan@example.com"},
		{Name: "María García", Email: "maria@example.com"},
		{Name: "Carlos López", Email: "carlos@example.com"},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}

	// Seed products
	products := []domain.Product{
		{Name: "Laptop Dell XPS 13", Price: 1200.00, Stock: 15},
		{Name: "iPhone 15 Pro", Price: 999.00, Stock: 25},
		{Name: "Sony WH-1000XM5", Price: 399.00, Stock: 30},
		{Name: "Samsung Galaxy Tab S9", Price: 649.00, Stock: 20},
		{Name: "Apple Watch Series 9", Price: 429.00, Stock: 40},
		{Name: "Logitech MX Master 3S", Price: 99.00, Stock: 50},
		{Name: "LG UltraFine 4K Monitor", Price: 699.00, Stock: 10},
		{Name: "Mechanical Keyboard RGB", Price: 159.00, Stock: 35},
	}
	if err := db.Create(&products).Error; err != nil {
		return err
	}

	log.Println("Database seeded successfully")
	return nil
}
