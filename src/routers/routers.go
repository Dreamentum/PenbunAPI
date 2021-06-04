package routers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"penbun.com/api/src/config"
	"penbun.com/api/src/model"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

// v1 version routing
func V1(engine *gin.Engine) {

	v1 := engine.Group("/v1")

	v1.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	v1.GET("/mssql", func(ctx *gin.Context) {
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
	})

	v1.GET("/getbooks", func(ctx *gin.Context) {
		rows, err := config.DB.Query("SELECT BookId, BookName, BookPrice, PublisherId, UpdateDate FROM tb_books ORDER BY BookId ASC")
		if err != nil {
			log.Fatalln("[SQL][Error][getAllBooks] SELECT failed:", err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var books []model.BOOK

		for rows.Next() {
			var b model.BOOK
			err := rows.Scan(&b.ID, &b.Name, &b.Price, &b.PublisherID, &b.UpdatedDate)

			if err != nil {
				log.Fatalln("[SQL][Error][getAllBooks] Scan failed:", err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
			books = append(books, b)
		}
		ctx.JSON(http.StatusOK, gin.H{"data": books})
		log.Println("[SQL][getAllBooks] SELECTED ALL BOOKS")
	})

	v1.GET("/getbookbyid", func(ctx *gin.Context) {
		var book model.BOOK
		book_id := ctx.Param("id")

		err := config.DB.QueryRow("SELECT BookId, BookName, BookPrice, PublisherId, UpdateDate FROM tb_books WHERE BookId=?", book_id).Scan(
			&book.ID, &book.Name, &book.Price, &book.PublisherID, &book.UpdatedDate,
		)

		if err != nil {
			log.Fatalln("[SQL][Error][getBookById] Setection Error :", err.Error())
		}
		ctx.JSON(http.StatusOK, gin.H{"book": book})
		log.Println("[SQL][getBookById] SELECTED BOOK: {} ", book_id)
	})

	v1.GET("/addgroup", func(ctx *gin.Context) {
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
	})

}

/**
Initialize routing, external call
*/
func Router(engine *gin.Engine) *gin.Engine {
	V1(engine)
	return engine
}
