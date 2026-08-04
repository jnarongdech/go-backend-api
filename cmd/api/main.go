package main

import (
	"log"

	"github.com/jnarongdech/go-backend-api/internal/config"
	"github.com/jnarongdech/go-backend-api/internal/handler"
	"github.com/jnarongdech/go-backend-api/internal/repository"
	"github.com/jnarongdech/go-backend-api/internal/service"
	"github.com/jnarongdech/go-backend-api/pkg/database"
	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 0. โหลดไฟล์ .env ทันทีที่เริ่มโปรแกรม
	config.LoadEnv()

	// ---------------------------------------------------------
	// 1. โหลดตัวแปรระบบและเชื่อมต่อ Database
	// ---------------------------------------------------------
	// config.LoadEnv() // สมมติว่ามีฟังก์ชันโหลดไฟล์ .env

	db := database.ConnectPostgres()
	defer db.Close() // สั่งปิดเมื่อโปรแกรมหยุดทำงาน

	// ---------------------------------------------------------
	// 2. เริ่มต้นใช้งาน sqlc
	// ---------------------------------------------------------
	// โยน *sql.DB เข้าไปในฟังก์ชัน New ที่ sqlc สร้างเตรียมไว้ให้
	sqlcQueries := repository.New(db)

	// ---------------------------------------------------------
	// 3. Dependency Injection (ประกอบร่าง)
	// ---------------------------------------------------------
	// ชั้น DB: โยน sqlcQueries ให้ Repository ใช้งาน
	userRepo := repository.NewUserRepository(sqlcQueries)

	// ชั้น Logic: โยน Repository ให้ Service ใช้งาน
	userService := service.NewUserService(userRepo)

	// ชั้น HTTP: โยน Service ให้ Handler ใช้งาน
	userHandler := handler.NewUserHandler(userService)

	// ---------------------------------------------------------
	// 4. ติดตั้ง Fiber Router
	// ---------------------------------------------------------

	// ---------------------------------------------------------
	// 🟢 สิ่งที่ต้องทำต่อเริ่มจากตรงนี้ครับ:
	// ---------------------------------------------------------

	// 1. สร้าง Fiber App
	app := fiber.New()

	// 2. ใส่ Middleware พื้นฐาน (Logger & CORS สำหรับให้ Next.js ยิงเข้าได้)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080", // เปลี่ยนเป็นโดเมน Frontend เมื่อ Deploy
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 3. กำหนด Routing
	api := app.Group("/api/v1")

	// Public Routes (ไม่ต้องต่อ Token)
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "message": "Server is healthy!"})
	})

	// 🔴 Protected Routes (ต้องผ่านการตรวจ Supabase JWT Token)
	// สมมติว่ามี Supabase Auth Middleware ในโฟลเดอร์ internal/middleware
	userRoutes := api.Group("/users" /*, middleware.SupabaseAuth() */)

	// นำ userHandler ไปผูกกับ HTTP Method และ Endpoint
	userRoutes.Get("/me", userHandler.GetProfile) // Endpoint ดึงข้อมูลตัวเอง
	userRoutes.Post("/", userHandler.CreateUser)  // Endpoint สร้าง User

	// 4. เริ่มรัน Server
	log.Println("Server is running on port:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
