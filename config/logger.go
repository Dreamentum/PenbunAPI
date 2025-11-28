package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger ตัวแปร Logger ใช้งานทั่วระบบ
var Logger *logrus.Logger

// InitLogger ฟังก์ชันสำหรับตั้งค่า Logger
func InitLogger() {
	Logger = logrus.New()

	// เปิดไฟล์ log สำหรับเขียนข้อมูล
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "logs/transaction.log" // ค่าเริ่มต้น
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Logger.Fatal("Failed to open log file: ", err)
	}

	// ตั้งค่า output เป็นไฟล์
	Logger.SetOutput(file)
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)
}
