package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"penbun.com/api/src/config"
	"penbun.com/api/src/controller"
	"penbun.com/api/src/middleware"
	"penbun.com/api/src/service"

	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
)

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	router := gin.New()
	router.Use(gin.Logger()) // Install the default logger, not required
	// router := gin.Default()
	// router.Use(cors.Default())
	router.Use(nice.Recovery(controller.RecoveryHandler))
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

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	router.GET("/panic", func(c *gin.Context) {
		// panic with a string -- the custom middleware could save this to a database or report it to the user
		panic("golang panic")
	})

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

	v1 := router.Group("/api/v1")

	v1.Use(middleware.AuthorizeJWT())
	{
		v1.StaticFile("/favicon.ico", "./favicon.ico")
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

	port := "443"
	router.Run(":" + port)
}
