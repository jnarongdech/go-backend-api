package server

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/jnarongdech/go-backend-api/internal/handler"
	"github.com/jnarongdech/go-backend-api/internal/repository"
	"github.com/jnarongdech/go-backend-api/internal/router"
	"github.com/jnarongdech/go-backend-api/internal/service"
)

func SetupServer(db *sql.DB) *fiber.App {
	sqlcQueries := repository.New(db)
	store := repository.NewStore(db)

	userRepo := repository.NewUserRepository(sqlcQueries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	productService := service.NewProductService(store)
	productHandler := handler.NewProductHandler(productService)

	customerService := service.NewCustomerService(store)
	customerHandler := handler.NewCustomerHandler(customerService)

	orderService := service.NewOrderService(store)
	orderHandler := handler.NewOrderHandler(orderService)

	handlers := &router.AppHandlers{
		User:     userHandler,
		Product:  productHandler,
		Customer: customerHandler,
		Order:    orderHandler,
	}

	app := fiber.New(fiber.Config{
		AppName: "STEEL-FACTORY API v1.0",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://127.0.0.1:3000", // Next.js or use * for dev
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Swagger Route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Mapping routing
	router.SetupRoutes(app, handlers)

	return app
}
