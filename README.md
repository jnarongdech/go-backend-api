# Steel Factory API (ชื่อโปรเจกต์ของคุณ)

A robust backend RESTful API built for managing inventory, products, and orders. Designed with scalability and clean code principles in mind.

## Key Features

- **Clean Architecture:** โครงสร้างโปรเจกต์แบ่ง Layer ชัดเจน (Handler, Service, Repository)
- **High Performance:** พัฒนาด้วย Go และ Fiber Framework ที่ประมวลผลได้รวดเร็ว
- **Type-Safe Database:** จัดการฐานข้อมูลอย่างปลอดภัยและลดข้อผิดพลาดด้วย `sqlc`
- **Transaction Management:** รองรับการทำ Database Transaction สำหรับระบบ Order อย่างรัดกุม
- **Frontend Ready:** โครงสร้าง JSON Response ถูกออกแบบมาให้เชื่อมต่อกับ Frontend ยุคใหม่อย่าง Next.js และ React ผ่าน TypeScript ได้อย่างลงตัว
- **API Documentation:** มีเอกสาร Swagger สำหรับทดสอบ API ครบทุก Endpoint

## Tech Stack

- **Language:** Go (Golang)
- **Framework:** Fiber
- **Database:** PostgreSQL (Supabase)
- **DB Toolkit:** sqlc
- **API Docs:** swaggo / Swagger

## Getting Started

### Prerequisites

- Go 1.20 หรือใหม่กว่า
- PostgreSQL Database (หรือ Supabase)

### Installation

1. Clone the repository:
   git clone https://github.com/jnarongdech/go-backend-api.git
   cd go-backend-api
2. ตั้งค่าตัวแปร Environment:
   สร้างไฟล์ .env ไว้ที่ root directory และเพิ่มข้อมูลเชื่อมต่อฐานข้อมูล
   DATABASE_URL=postgres://user:password@host:port/dbname
3. ติดตั้ง Dependencies และรันเซิร์ฟเวอร์:
   go mod tidy
   swag init -g cmd/api/main.go -d . --parseInternal
   go run cmd/api/main.go
   เซิร์ฟเวอร์จะรันที่พอร์ต
   > > http://localhost:8080

### API Documentation

เมื่อรันเซิร์ฟเวอร์เรียบร้อยแล้ว สามารถเข้าดูเอกสารและทดสอบ API ได้ที่:

> > http://localhost:8080/swagger/index.html

### Project Structure

- /cmd - จุดเริ่มต้นการทำงานของแอปพลิเคชัน (main.go)
- /internal/handler - รับ HTTP Request และตอบกลับเป็น JSON
- /internal/service - จัดการ Business Logic หลักของระบบ
- /internal/repository - จัดการการดึง/บันทึกข้อมูลกับ Database (sqlc)
- /internal/dto - โครงสร้างข้อมูล (Data Transfer Object)
- /docs - ไฟล์ที่ถูก Generate โดย Swagger

### Author

[Narongdech Petchtra]

- Github: jnarongdech
- LinkedIn: www.linkedin.com/in/jnarongdech78
