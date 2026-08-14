package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jnarongdech/go-backend-api/internal/handler"
)

type AppHandlers struct {
	User     *handler.UserHandler
	Product  *handler.ProductHandler
	Customer *handler.CustomerHandler
	Order    *handler.OrderHandler
}

func SetupRoutes(app *fiber.App, h *AppHandlers) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "STEEL-FACTORY API is healthy!",
		})
	})

	v1 := app.Group("/api/v1")

	users := v1.Group("/users")
	users.Post("/", h.User.CreateUser)
	users.Patch("/", h.User.UpdateUser)
	// users.Get("/", h.User.GetAllUsers)
	users.Delete("/", h.User.SoftDeleteUser)

	products := v1.Group("/products")
	products.Get("/", h.Product.GetProducts)
	products.Post("/", h.Product.CreateProduct)
	// products.Get("/:id", h.Product.GetProductByID)

	customers := v1.Group("/customers")
	customers.Get("/:id", h.Customer.GetCustomerByID)
	customers.Post("/", h.Customer.CreateCustomer)
	customers.Put("/:id", h.Customer.UpdateCustomer)

	orders := v1.Group("/orders")
	orders.Get("/", h.Order.GetOrders)
	orders.Get("/:id", h.Order.GetOrderByID)
	orders.Post("/", h.Order.CreateOrderWithItems)

	// ---------------------------------------------------------
	// ตัวอย่าง: ถ้าอนาคตมีระบบที่ต้อง Login ถึงจะเข้าได้
	// ---------------------------------------------------------
	// secure := v1.Group("/secure", middleware.RequireAuth())
	// secure.Get("/dashboard", h.User.GetDashboardStats)
}
