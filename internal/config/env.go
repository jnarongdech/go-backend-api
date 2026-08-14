package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// พยายามโหลดไฟล์ .env (สำหรับรันในเครื่อง)
	// แต่ถ้าโหลดไม่ได้ (เช่นรันบน Render) ให้ข้ามเลย
	_ = godotenv.Load()

	// เช็คค่าจริงๆ ว่ามีตัวแปรถูกส่งมาไหม (ไม่ว่าจะจาก .env หรือจาก Render)
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// ถ้าตรงนี้ไม่มีค่า ค่อยโวยวายของจริงครับ เพราะแปลว่าไม่ได้ตั้งค่าไว้เลย
		log.Fatal("ERROR: DATABASE_URL is not set in environment")
	}
}
