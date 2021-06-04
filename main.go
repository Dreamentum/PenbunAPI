package main

import (
	"database/sql"
	"log"
	"net/http"

	"penbun.com/api/src/config"
	"penbun.com/api/src/controller"
	"penbun.com/api/src/middleware"
	"penbun.com/api/src/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// router := gin.New()
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	var loginService service.LoginService = service.StaticLoginService()
	var jwtService service.JWTService = service.JWTAuthService()
	var loginController controller.LoginController = controller.LoginHandler(loginService, jwtService)
	var err error

	config.DB, err = sql.Open("mssql", config.ConnectionString)
	if err != nil {
		log.Fatalln("[SQL][Error] Open connection failed:", err.Error())
	}

	log.Printf("[MSSQL] Connected!\n")
	defer config.DB.Close()

	router.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	v1 := router.Group("/v1/api")

	v1.Use(middleware.AuthorizeJWT())
	{
		v1.GET("/welcome", controller.HelloHandler)
		v1.GET("/mssql", controller.CheckMssql)
		v1.GET("/book/select", controller.GetBooks)
		v1.GET("/book/:id", controller.GetBookById)
		v1.GET("/book/type/select", controller.GetBookTypes)
		v1.POST("/book/type/add", controller.AddBookType)
		v1.POST("/book/type/delete", controller.DeleteBookType)
		v1.POST("/book/type/update", controller.UpdateBookType)
		v1.POST("/group/add", controller.AddGroup)
		v1.GET("/publisher/select", controller.GetPublishers)
	}

	port := "8080"
	router.Run(":" + port)
}
