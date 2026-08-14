package main

import (
	"log"
	"os" // 🟢 นำเข้า os เพื่อใช้อ่านค่า PORT จาก env

	"github.com/jnarongdech/go-backend-api/internal/config"
	"github.com/jnarongdech/go-backend-api/internal/server"
	"github.com/jnarongdech/go-backend-api/pkg/database"
	_ "github.com/lib/pq"

	// swagger

	_ "github.com/jnarongdech/go-backend-api/docs"
)

// @title GO-BACKEND-API
// @version 1.0
// @description ระบบจัดการ API สำหรับโรงงาน STEEL-FACTORY
// @host localhost:8080
// @BasePath /
func main() {
	config.LoadEnv()

	db := database.ConnectPostgres()
	defer db.Close()

	app := server.SetupServer(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
