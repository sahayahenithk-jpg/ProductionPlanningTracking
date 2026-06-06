package main

import (
	"log"

	"projectplanningtracking/backend/db"
	"projectplanningtracking/backend/handlers"
	"projectplanningtracking/backend/middleware"
	"projectplanningtracking/backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	database := db.ConnectDB()
	database.AutoMigrate(&models.User{}, &models.Product{}, &models.ProductionPlan{}, &models.ProductionEntry{})

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/api/auth/register", handlers.Register(database))
	r.POST("/api/auth/login", handlers.Login(database))
	// also expose simpler endpoints
	r.POST("/api/register", handlers.Register(database))
	r.POST("/api/login", handlers.Login(database))

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/user", handlers.Profile(database))

	// User management (admin only)
	protected.GET("/users", middleware.RequireRole("admin"), handlers.ListUsers(database))
	protected.GET("/users/:id", middleware.RequireRole("admin"), handlers.GetUser(database))
	protected.POST("/users", middleware.RequireRole("admin"), handlers.CreateUser(database))
	protected.PUT("/users/:id", middleware.RequireRole("admin"), handlers.UpdateUser(database))
	protected.DELETE("/users/:id", middleware.RequireRole("admin"), handlers.DeleteUser(database))

	// Products CRUD
	protected.GET("/products", handlers.ListProducts(database))
	protected.GET("/products/:id", handlers.GetProduct(database))
	protected.POST("/products", middleware.RequireRole("admin"), handlers.CreateProduct(database))
	protected.PUT("/products/:id", middleware.RequireRole("admin"), handlers.UpdateProduct(database))
	protected.DELETE("/products/:id", middleware.RequireRole("admin"), handlers.DeleteProduct(database))

	protected.GET("/plans", handlers.ListPlans(database))
	protected.POST("/plans", middleware.RequireRole("planner", "admin"), handlers.CreatePlan(database))
	protected.PUT("/plans/:id", middleware.RequireRole("planner", "admin"), handlers.UpdatePlan(database))

	protected.GET("/production", handlers.ListProductionEntries(database))
	protected.POST("/production", middleware.RequireRole("operator", "admin"), handlers.CreateProductionEntry(database))
	protected.PUT("/production/:id", middleware.RequireRole("operator", "admin"), handlers.UpdateProductionEntry(database))

	protected.GET("/reports/variance", middleware.RequireRole("planner", "admin"), handlers.VarianceReport(database))
	protected.GET("/reports/summary", middleware.RequireRole("planner", "admin"), handlers.SummaryReport(database))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
