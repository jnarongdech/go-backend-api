package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// สั่งให้อ่านไฟล์ .env ที่อยู่หน้าสุดของโปรเจกต์
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or error loading it.")
	} else {
		log.Println("Environment variables loaded successfully!")
	}
}
