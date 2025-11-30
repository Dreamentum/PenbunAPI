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

// Register - ลงทะเบียนผู้ใช้ใหม่
func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		config.Logger.WithError(err).Warn("Register attempt: Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ตรวจสอบ username ว่ามีอยู่แล้วหรือไม่
	var exists bool
	// NOTE: การใช้ sql.Named ขึ้นอยู่กับ Driver ที่ใช้
	err := config.DB.QueryRow("SELECT 1 FROM tb_users WHERE user_name = @UserName",
		sql.Named("UserName", user.UserName)).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		config.Logger.WithError(err).Error("Database error during Register check")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	if exists {
		config.Logger.WithField("user_name", user.UserName).Warn("Register attempt: Username already exists")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Hash รหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		config.Logger.WithError(err).Error("Failed to hash password during Register")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// บันทึกข้อมูลลงฐานข้อมูล
	// แก้ไข: เพิ่ม UpdateDate ใน Exec
	_, err = config.DB.Exec("INSERT INTO tb_users (user_name, user_password, update_date) VALUES (@UserName, @UserPassword, @UpdateDate)",
		sql.Named("UserName", user.UserName),
		sql.Named("UserPassword", string(hashedPassword)),
		sql.Named("UpdateDate", time.Now()),
	)
	if err != nil {
		config.Logger.WithError(err).WithField("user_name", user.UserName).Error("Failed to register user to database")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	config.Logger.WithField("user_name", user.UserName).Info("User registered successfully")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// RefreshToken - สร้าง JWT Token ใหม่จาก Token เก่าที่ยังไม่หมดอายุ
func RefreshToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		config.Logger.Warn("RefreshToken attempt failed: Missing token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// แก้ไข: ใช้ TrimPrefix
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	tokenPrefix := ""
	if len(tokenString) > 10 {
		tokenPrefix = tokenString[:10] + "..."
	}

	if config.IsBlacklisted(tokenString) {
		config.Logger.WithField("token_prefix", tokenPrefix).Warn("Refresh attempt: Token is blacklisted")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is blacklisted"})
	}

	// ตรวจสอบความถูกต้องของ Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		config.Logger.WithError(err).WithField("token_prefix", tokenPrefix).Warn("Refresh attempt: Invalid or expired token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// ดึง Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		config.Logger.Error("Failed to extract claims from token during Refresh")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// สร้าง Token ใหม่
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": claims["user_name"],
		"iss":       claims["iss"],
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err = newToken.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
		config.Logger.WithError(err).Error("Failed to generate new token during Refresh")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	config.Logger.WithField("user_name", claims["user_name"]).Info("Token refreshed successfully")

	return c.JSON(fiber.Map{"token": tokenString})
}

// Login - เข้าสู่ระบบและสร้าง JWT Token
func Login(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		config.Logger.WithError(err).Warn("Login attempt failed: Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid request",
		})
	}

	config.Logger.WithField("user_name", user.UserName).Info("Login attempt received")
	log.Println("[INFO] Login attempt received for user:", user.UserName)

	// ตรวจสอบ username และ password จาก database
	var hashedPassword string
	err := config.DB.QueryRow("SELECT user_password FROM tb_users WHERE user_name = @UserName",
		sql.Named("UserName", user.UserName)).Scan(&hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			config.Logger.WithField("user_name", user.UserName).Warn("Login attempt: Username not found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail",
				"error":  "Invalid username or password",
			})
		}
		config.Logger.WithError(err).Error("Database error during Login query")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Database error",
		})
	}

	// ตรวจสอบรหัสผ่าน
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			config.Logger.WithField("user_name", user.UserName).Warn("Login attempt: Password mismatch")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail",
				"error":  "Invalid username or password",
			})
		}
		config.Logger.WithError(err).WithField("user_name", user.UserName).Error("Error verifying password")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Internal error during password check",
		})
	}

	config.Logger.WithField("user_name", user.UserName).Info("Password verified successfully")
	log.Println("[INFO] Password verified successfully for user:", user.UserName)

	// สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": user.UserName,
		"iss":       "penbun-api",
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
	})

	// Signed token
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
		config.Logger.WithError(err).Error("Failed to sign token during Login")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to create token",
		}) // <--- แก้ไข Syntax Error ตรงนี้
	}

	// Logrus: บันทึกเมื่อ Login สำเร็จ
	config.Logger.WithField("user_name", user.UserName).Info("User logged in successfully")

	// ส่ง Token กลับไป
	return c.JSON(fiber.Map{
		"status": "success",
		"token":  tokenString,
	})
} // <--- ต้องปิด func Login ที่นี่

// Logout - เพิ่ม Token ลงใน Blacklist
func Logout(c *fiber.Ctx) error { // <--- เริ่ม func Logout ที่ถูกต้อง

	// รับ Token จาก Header Authorization
	token := c.Get("Authorization")
	if token == "" {
		config.Logger.Warn("Logout attempt failed: Missing Authorization token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing token",
		})
	}

	// ตัดคำว่า "Bearer " ออกจาก Token
	token = strings.TrimPrefix(token, "Bearer ")

	// ตัด Token ให้สั้นลงเพื่อความปลอดภัยในการ Log
	tokenPrefix := token
	if len(token) > 10 {
		tokenPrefix = token[:10] + "..."
	}

	// เพิ่ม Token ลงใน Blacklist
	if !config.IsBlacklisted(token) {
		// สมมติว่าฟังก์ชันนี้มีอยู่จริงใน config
		config.AddToBlacklist(token) 

		// Logrus: บันทึกเมื่อ Token ถูกเพิ่มลง Blacklist สำเร็จ
		config.Logger.WithFields(logrus.Fields{
			"action":       "blacklist_add",
			"token_prefix": tokenPrefix,
		}).Info("User token blacklisted successfully (Logout)")
	} else {
		// Logrus: บันทึกเมื่อพยายาม Logout ด้วย Token ที่ถูก Blacklist แล้ว
		config.Logger.WithFields(logrus.Fields{
			"action":       "blacklist_duplicate",
			"token_prefix": tokenPrefix,
		}).Warn("Logout attempt with an already blacklisted token")
	}

	// Logrus: บันทึกเมื่อ Logout สำเร็จ
	config.Logger.Info("Logout process completed")
	log.Println("[INFO] Logout successful")

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}