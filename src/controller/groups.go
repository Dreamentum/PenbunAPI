package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"penbun.com/api/src/config"
)

func AddGroup(ctx *gin.Context) {
	prefix := ctx.Request.FormValue("prefix")
	book_group_name := ctx.Request.FormValue("groupname")
	book_group_level := ctx.Request.FormValue("grouplevel")
	book_group_description := ctx.Request.FormValue("description")
	book_group_status := ctx.Request.FormValue("status")
	book_group_create_by := ctx.Request.FormValue("createby")
	book_group_update_by := ctx.Request.FormValue("updateby")

	// EXEC dbo.sp_GroupInsert 'GBK','TEST', '1', 'Test how to call stored procedure', 1, 'ADMIN', 'ADMIN'

	_, err := config.DB.Query("EXEC sp_GroupInsert ?, ?, ?, ?, ?, ?, ?",
		prefix, book_group_name, book_group_level, book_group_description,
		book_group_status, book_group_create_by, book_group_update_by)

	if err != nil {
		log.Fatalln("[!] Insert group failed:", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"DONE": book_group_name})
	log.Println("[+] GROUP NAME: ", book_group_name)
}
