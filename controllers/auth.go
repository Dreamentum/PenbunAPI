package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"

	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus" // Logrus สำหรับบันทึก Log

	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// ตรวจสอบ username ว่ามีอยู่แล้วหรือไม่
	var exists bool
	err := config.DB.QueryRow("SELECT 1 FROM tb_users WHERE user_name = @UserName",
		sql.Named("UserName", user.UserName)).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Hash รหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// บันทึกข้อมูลลงฐานข้อมูล
	_, err = config.DB.Exec("INSERT INTO tb_users (user_name, user_password, update_date) VALUES (@UserName, @UserPassword)",
		sql.Named("UserName", user.UserName),
		sql.Named("UserPassword", string(hashedPassword)),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

/*
ฟังก์ชัน RefreshToken ใช้สำหรับออก JWT Token ใหม่
เมื่อ Token ปัจจุบันใกล้หมดอายุ โดยยังรักษาข้อมูลผู้ใช้งานเดิมไว้
วิธีการทำงานคือ ตรวจสอบ Token ปัจจุบันก่อน
จากนั้นสร้าง Token ใหม่ที่มีวันหมดอายุยาวขึ้น.
*/
func RefreshToken(c *fiber.Ctx) error {
	// รับ Token จาก Header Authorization
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// ตัดคำว่า "Bearer " ออกจาก Token
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// ตรวจสอบว่า Token อยู่ใน Blacklist หรือไม่
	if config.IsBlacklisted(tokenString) {
		log.Println("[DEBUG] Token is blacklisted:", tokenString)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is blacklisted"})
	}

	// ตรวจสอบความถูกต้องของ Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// สร้าง Token ใหม่
	claims := token.Claims.(jwt.MapClaims)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": claims["user_name"], // คงค่าเดิมจาก Token เก่า
		"iss":       claims["iss"],
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err = newToken.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

func Login(c *fiber.Ctx) error {
	// รับ username และ password จาก request body
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid request",
		})
	}

	// เพิ่ม Log สำหรับ Debug
	log.Println("[INFO] Username from Request:", user.UserName)
	log.Println("[INFO] Password from Request:", user.Password)

	// ตรวจสอบ username และ password จาก database
	var hashedPassword string
	err := config.DB.QueryRow("SELECT user_password FROM tb_users WHERE user_name = @UserName",
		sql.Named("UserName", user.UserName)).Scan(&hashedPassword)

	if err != nil {
		// กรณีไม่มีข้อมูลในฐานข้อมูล
		if err == sql.ErrNoRows {
			log.Println("[SQL] Username not found:", user.UserName)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail",
				"error":  "Invalid username or password",
			})
		}
		// กรณี Error อื่นๆ
		log.Println("[SQL] Error querying database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Database error",
		})
	}

	// ตรวจสอบรหัสผ่าน
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		// กรณีรหัสผ่านไม่ตรงกัน
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("[DEBUG] Password mismatch for user:", user.UserName)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail",
				"error":  "Invalid Password",
			})
		}

		// กรณีเกิดข้อผิดพลาดอื่นๆ ระหว่างการตรวจสอบรหัสผ่าน
		log.Println("[DEBUG] Error verifying password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Database error",
		})
	}
	// เพิ่ม Log กรณีตรวจสอบสำเร็จ
	log.Println("[INFO] Password verified successfully for user:", user.UserName)

	// สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": user.UserName,
		"iss":       "PenbunAPI",                          // ชื่อระบบที่ออกโทเค็น
		"iat":       time.Now().Unix(),                    // เวลาที่ออกโทเค็น
		"exp":       time.Now().Add(time.Hour * 1).Unix(), // เวลาหมดอายุของโทเค็น
	})
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Could not generate token",
		})
	}

	// เพิ่ม Log กรณีสร้าง JWT สำเร็จ
	log.Println("[INFO] Generate JWT successfully for user:", user.UserName)
	config.Logger.WithFields(logrus.Fields{
		"user_name": user.UserName,
		"iss":       "PenbunAPI",                          // ชื่อระบบที่ออกโทเค็น
		"iat":       time.Now().Unix(),                    // เวลาที่ออกโทเค็น
		"exp":       time.Now().Add(time.Hour * 1).Unix(), // เวลาหมดอายุของโทเค็น
	}).Info("JWT created successfully")

	// ส่ง JWT token กลับไป
	return c.JSON(fiber.Map{
		"status":  "success",
		"token":   tokenString,
		"message": "Login successful",
	})
}

func Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing token",
		})
	}

	// ตัดคำว่า "Bearer " ออกจาก Token
	token = strings.TrimPrefix(token, "Bearer ")

	// เพิ่ม Token ลงใน Blacklist
	if !config.IsBlacklisted(token) {
		config.AddToBlacklist(token)
	}

	// Log ว่า Token ถูกเพิ่มแล้ว
	log.Println("[DEBUG] Token added to blacklist:", token)
	log.Println("[INFO] Logout successfully for user:", user.UserName)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}
