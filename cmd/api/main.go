package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"ecommerce/internal/database"
	"ecommerce/internal/order"
	"ecommerce/internal/payment"
	"ecommerce/internal/product"
	"ecommerce/internal/user"
)

func main() {

	db, err := database.OpenDB("ecommerce.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo := user.NewSQLiteUserRepo(db)
	productRepo := product.NewSQLiteProductRepo(db)
	orderRepo := order.NewSQLiteOrderRepo(db)
	paymentRepo := payment.NewSQLitePaymentRepo(db)

	userService := user.NewService(userRepo)
	productService := product.NewService(productRepo)
	orderService := order.NewService(orderRepo)
	paymentService := payment.NewService(paymentRepo)

	userHandler := user.NewHandler(userService)
	productHandler := product.NewHandler(productService)
	orderHandler := order.NewHandler(orderService)
	paymentHandler := payment.NewHandler(paymentService)

	r := gin.Default()

	r.StaticFS("/ui", http.Dir("./ui"))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/index.html")
	})

	r.POST("/api/v1/users", userHandler.CreateUser)
	r.GET("/api/v1/users/:id", userHandler.GetUser)
	r.PUT("/api/v1/users/:id", userHandler.UpdateUser)
	r.DELETE("/api/v1/users/:id", userHandler.DeleteUser)

	r.POST("/api/v1/products", productHandler.CreateProduct)
	r.GET("/api/v1/products/:id", productHandler.GetProduct)
	r.PUT("/api/v1/products/:id", productHandler.UpdateProduct)
	r.DELETE("/api/v1/products/:id", productHandler.DeleteProduct)
	r.GET("/api/v1/products", productHandler.ListProducts)

	r.POST("/api/v1/orders", orderHandler.CreateOrder)
	r.GET("/api/v1/orders/:id", orderHandler.GetOrder)
	r.PUT("/api/v1/orders/:id/status", orderHandler.UpdateOrderStatus)
	r.DELETE("/api/v1/orders/:id", orderHandler.DeleteOrder)

	r.GET("/api/v1/orders", orderHandler.ListOrders)

	r.GET("/api/v1/users/:id/orders", orderHandler.GetOrdersByUser)

	r.POST("/api/v1/payments", paymentHandler.CreatePayment)
	r.GET("/api/v1/payments/:id", paymentHandler.GetPayment)
	r.PUT("/api/v1/payments/:id/confirm", paymentHandler.ConfirmPayment)
	r.PUT("/api/v1/payments/:id/fail", paymentHandler.FailPayment)
	r.GET("/api/v1/users/:id/payments", paymentHandler.GetPaymentsByUser)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
