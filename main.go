package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"PenbunAPI/config"
	"PenbunAPI/middleware"
	"PenbunAPI/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/joho/godotenv"
)

func init() {
	// โหลดไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}
}

func cleanUp() {
	fmt.Println("Cleaning Up..")
	for {
		// deletes old files here
		time.Sleep(60 * time.Second)
	}
}

func main() {
	// 4 procs/childs max
	runtime.GOMAXPROCS(3)

	// start a cleanup cron-job
	go cleanUp()

	// เริ่มต้น Logger
	config.InitLogger()

	// อ่าน PORT จากไฟล์ .env
	port := os.Getenv("FIBER_PORT")
	if port == "" {
		port = "8089" // default if no port in .env
	}

	// สร้าง Fiber App
	app := fiber.New(fiber.Config{
		Prefork:           false,
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: true,
		ServerHeader:      "Fiber",
		AppName:           "PENBUN API v.1.9.5",
	})

	// ✅ Serve favicon.ico
	// app.Get("/favicon.ico", func(c *fiber.Ctx) error {
	// 	return c.SendStatus(fiber.StatusNoContent) // หรือใช้ StatusOK ก็ได้
	// })

	// ใช้ JWTMiddleware ระดับ Global
	// app.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))
	// log.Println("[DEBUG] JWT is :", token)

	// เพิ่ม Logger Middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use(middleware.NewLoggerMiddleware())

	// เชื่อมต่อ Database
	config.ConnectDatabase()

	// ลงทะเบียน Routes พร้อมส่ง Database Connection
	routes.RegisterV1Routes(app, config.DB)

	// routes.RegisterV2Routes(app, config.DB)

	// เริ่มเซิร์ฟเวอร์
	log.Println("Starting server on port", port)
	log.Fatal(app.Listen(":" + port))

	// จัดการ Graceful Shutdown
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Error starting server: %v\n", err)
		}
	}()

	// รอรับสัญญาณ Interrupt หรือ Kill
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	fmt.Println("Gracefully shutting down...")

	// รอให้ Fiber Shutdown อย่างปลอดภัย
	if err := app.Shutdown(); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}

	// ปิดฐานข้อมูลหรือกระบวนการที่ค้างอยู่
	config.DB.Close()
	fmt.Println("Cleanup completed.")
}
