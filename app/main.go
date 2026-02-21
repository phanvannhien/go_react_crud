package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nhienphan/full_rest_app/db/sqlc"
	"github.com/nhienphan/full_rest_app/internal/handler"
	"github.com/nhienphan/full_rest_app/internal/middleware"
	"github.com/nhienphan/full_rest_app/internal/routes"
	"github.com/nhienphan/full_rest_app/internal/service"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	// Required env vars
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to PostgreSQL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Initialize sqlc queries
	queries := sqlc.New(pool)

	// Initialize services (dependency injection)
	authService := service.NewAuthService(queries, jwtSecret)
	userService := service.NewUserService(queries)
	productService := service.NewProductService(queries)
	orderService := service.NewOrderService(queries, pool)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	// Create Echo instance
	e := echo.New()

	// Global middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Request body size limit (1MB)
	e.Use(echoMiddleware.BodyLimit("1M"))

	// Register public routes (no auth)
	routes.RegisterAuthRoutes(e, authHandler)

	// Protected routes group
	api := e.Group("/api", middleware.JWTAuth(jwtSecret))
	routes.RegisterUserRoutes(api, userHandler)
	routes.RegisterProductRoutes(api, productHandler)
	routes.RegisterOrderRoutes(api, orderHandler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
