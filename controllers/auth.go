package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"

	"database/sql"
	// ลบ "log" ออก หรือไม่ใช้มันเลย เพื่อรักษามาตรฐาน Logrus 
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
        config.Logger.WithError(err).Warn("Register attempt: Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// ตรวจสอบ username ว่ามีอยู่แล้วหรือไม่
	var exists bool
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
	_, err = config.DB.Exec("INSERT INTO tb_users (user_name, user_password, update_date) VALUES (@UserName, @UserPassword)",
		sql.Named("UserName", user.UserName),
		sql.Named("UserPassword", string(hashedPassword)),
	)
	if err != nil {
        config.Logger.WithError(err).WithField("user_name", user.UserName).Error("Failed to register user to database")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

    // Logrus: บันทึกเมื่อลงทะเบียนสำเร็จ
    config.Logger.WithField("user_name", user.UserName).Info("User registered successfully")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

func RefreshToken(c *fiber.Ctx) error {
	// รับ Token จาก Header Authorization
	tokenString := c.Get("Authorization")
	if tokenString == "" {
        config.Logger.Warn("RefreshToken attempt failed: Missing token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// ตัดคำว่า "Bearer " ออกจาก Token
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

    tokenPrefix := ""
    if len(tokenString) > 10 {
        tokenPrefix = tokenString[:10] + "..." 
    }

	// ตรวจสอบว่า Token อยู่ใน Blacklist หรือไม่
	if config.IsBlacklisted(tokenString) {
        // Logrus: บันทึกเมื่อ Token อยู่ใน Blacklist
        config.Logger.WithField("token_prefix", tokenPrefix).Warn("Refresh attempt: Token is blacklisted")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is blacklisted"})
	}

	// ตรวจสอบความถูกต้องของ Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
        // Logrus: บันทึกเมื่อ Token ไม่ถูกต้อง
        config.Logger.WithError(err).WithField("token_prefix", tokenPrefix).Warn("Refresh attempt: Invalid or expired token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// สร้าง Token ใหม่
	claims := token.Claims.(jwt.MapClaims)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": claims["user_name"], // คงค่าเดิมจาก Token เก่า
		"iss": 		 claims["iss"],
		"iat": 		 time.Now().Unix(),
		"exp": 		 time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err = newToken.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
        config.Logger.WithError(err).Error("Failed to generate new token during Refresh")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}
    
    // Logrus: บันทึกเมื่อ Refresh Token สำเร็จ
    config.Logger.WithField("user_name", claims["user_name"]).Info("Token refreshed successfully")

	return c.JSON(fiber.Map{"token": tokenString})
}

func Login(c *fiber.Ctx) error {
	// รับ username และ password จาก request body
	var user models.User
	if err := c.BodyParser(&user); err != nil {
        config.Logger.WithError(err).Warn("Login attempt failed: Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid request",
		})
	}
    
    // Logrus: บันทึกการรับคำขอ Login 
    config.Logger.WithField("user_name", user.UserName).Info("Login attempt received")

	// ตรวจสอบ username และ password จาก database
	var hashedPassword string
	err := config.DB.QueryRow("SELECT user_password FROM tb_users WHERE user_name = @UserName",
		sql.Named("UserName", user.UserName)).Scan(&hashedPassword)

	if err != nil {
		// กรณีไม่มีข้อมูลในฐานข้อมูล
		if err == sql.ErrNoRows {
            // Logrus: บันทึกเมื่อไม่พบผู้ใช้
            config.Logger.WithField("user_name", user.UserName).Warn("Login attempt: Username not found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail",
				"error":  "Invalid username or password",
			})
		}
		// กรณี Error อื่นๆ
        config.Logger.WithError(err).Error("Database error during Login query")
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
            // Logrus: บันทึกเมื่อรหัสผ่านไม่ตรง
            config.Logger.WithField("user_name", user.UserName).Warn("Login attempt: Password mismatch")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail",
				"error":  "Invalid Password",
			})
		}

		// กรณีเกิดข้อผิดพลาดอื่นๆ ระหว่างการตรวจสอบรหัสผ่าน
        config.Logger.WithError(err).WithField("user_name", user.UserName).Error("Error verifying password")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Internal error during password check",
		})
	}
    
    // Logrus: บันทึกเมื่อตรวจสอบรหัสผ่านสำเร็จ
    config.Logger.WithField("user_name", user.UserName).Info("Password verified successfully")

	// สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": user.UserName,
		"iss": 		 "PenbunAPI", 					// ชื่อระบบที่ออกโทเค็น
		"iat": 		 time.Now().Unix(), 			// เวลาที่ออกโทเค็น
		"exp": 		 time.Now().Add(time.Hour * 1).Unix(), // เวลาหมดอายุของโทเค็น
	})
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
        config.Logger.WithError(err).Error("Failed to generate JWT token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Could not generate token",
		})
	}

    // Logrus: บันทึกเมื่อสร้าง JWT สำเร็จ
	config.Logger.WithFields(logrus.Fields{
		"user_name": user.UserName,
		"iss": 		 "PenbunAPI",
		"exp": 		 time.Now().Add(time.Hour * 1).Unix(),
	}).Info("User logged in and JWT created successfully")

	// ส่ง JWT token กลับไป
	return c.JSON(fiber.Map{
		"status":  "success",
		"token": 	tokenString,
		"message": "Login successful",
	})
}

func Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		// Logrus: แจ้งเตือนเมื่อมีการเรียก Logout โดยไม่มี Token
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
		config.AddToBlacklist(token)
		
		// Logrus: บันทึกเมื่อ Token ถูกเพิ่มลง Blacklist สำเร็จ
		config.Logger.WithFields(logrus.Fields{
			"action": 	  "blacklist_add",
			"token_prefix": tokenPrefix,
		}).Info("User token blacklisted successfully (Logout)")
	} else {
		// Logrus: บันทึกเมื่อพยายาม Logout ด้วย Token ที่ถูก Blacklist แล้ว
		config.Logger.WithFields(logrus.Fields{
			"action": 	  "blacklist_duplicate",
			"token_prefix": tokenPrefix,
		}).Warn("Logout attempt with an already blacklisted token")
	}

	// Logrus: บันทึกเมื่อ Logout สำเร็จ (ไม่ว่า Token จะถูก Blacklist ซ้ำหรือไม่ก็ตาม)
	config.Logger.Info("Logout process completed")
    // ลบ log.Println(...) ที่ซ้ำซ้อนออก

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}