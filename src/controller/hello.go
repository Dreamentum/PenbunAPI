package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"penbun.com/api/src/config"
)

func HelloHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "PENBUN API Version 1.0.14.2 [2022/FEB/14]"})
}

func CheckMssql(ctx *gin.Context) {
	err := config.DB.PingContext(ctx)
	if err != nil {
		log.Fatalln("[SQL][Error][checkMssqlVersion] Ping database server failed:", err.Error())
	}

	var version string

	err = config.DB.QueryRowContext(ctx, "SELECT @@version").Scan(&version)
	if err != nil {
		log.Fatalln("[SQL][Error][checkMssqlVersion] Scan failed:", err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"message": version})
	log.Printf("[SQL] version %s\n", version)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func RecoveryHandler(c *gin.Context, err interface{}) {
	c.HTML(500, "error.tmpl", gin.H{
		"title": "Error",
		"err":   err,
	})
}
