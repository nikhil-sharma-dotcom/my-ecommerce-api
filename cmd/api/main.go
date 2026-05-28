package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/adapters/postgresql/sqlc"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/config"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/handlers"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/logger"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/middleware"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/orders"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/products"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/users"
)

func main() {
	cfg := config.Load()
	logger.Init()
	defer logger.Sync()

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, cfg.DBConnectionString())
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		logger.Log.Fatal("Failed to ping database", zap.Error(err))
	}

	queries := sqlc.New(dbpool)

	userService := users.NewService(queries, cfg)
	productService := products.NewService(queries)
	orderService := orders.NewService(queries)

config, err := pgxpool.ParseConfig(cfg.DBConnectionString())
if err != nil {
    logger.Log.Fatal("Failed to parse database config", zap.Error(err))
}
dbpool, err := pgxpool.NewWithConfig(ctx, config)	authHandler := handlers.NewAuthHandler(userService, cfg)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Logger)
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware)

	rateLimiter := middleware.NewIPRateLimiter(10, 20)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg))

			r.Get("/profile", authHandler.GetProfile)

			r.Get("/products", productHandler.GetAll)
			r.Get("/products/{id}", productHandler.GetByID)

			r.Post("/orders", orderHandler.Create)
			r.Get("/orders", orderHandler.GetMyOrders)
			r.Get("/orders/{id}", orderHandler.GetByID)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg))
			r.Use(middleware.RoleMiddleware("admin"))

			r.Post("/products", productHandler.Create)
			r.Put("/products/{id}", productHandler.Update)
			r.Delete("/products/{id}", productHandler.Delete)

			r.Patch("/orders/{id}/status", orderHandler.UpdateStatus)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Log.Info("Server starting", zap.String("address", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
	}
}
