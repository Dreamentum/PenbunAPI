package controller

import (
	"log"
	"net/http"

	"penbun.com/api/src/config"
	"penbun.com/api/src/model"

	"github.com/gin-gonic/gin"
)

// Book
func GetBooks(ctx *gin.Context) {
	rows, err := config.DB.Query("SELECT BookId, BookName, BookPrice, PublisherId, UpdateDate FROM tb_books ORDER BY BookId ASC")
	if err != nil {
		log.Fatalln("[SQL][Error][GetBooks] SELECT failed:", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []model.BOOK

	for rows.Next() {
		var b model.BOOK
		err := rows.Scan(&b.ID, &b.Name, &b.Price, &b.PublisherID, &b.UpdatedDate)
		if err != nil {
			log.Fatalln("[SQL][Error][GetBooks] Scan failed:", err.Error())
			return
		}
		books = append(books, b)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": books})
	log.Println("[SQL][GetBooks] SELECTED ALL BOOKS")
}

func GetBookById(ctx *gin.Context) {
	var book model.BOOK
	book_id := ctx.Param("id")

	err := config.DB.QueryRow("SELECT BookId, BookName, BookPrice, PublisherId, UpdateDate FROM tb_books WHERE BookId=?", book_id).Scan(
		&book.ID, &book.Name, &book.Price, &book.PublisherID, &book.UpdatedDate,
	)

	if err != nil {
		log.Fatalln("[SQL][Error][GetBookById] Setection Error :", err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"book": book})
	log.Println("[SQL][GetBookById] SELECTED BOOK: ", book_id)
}

// func GetBookByIdAndPrice(ctx *gin.Context) {
// 	var book model.BOOK
// 	book_id := ctx.Param("id")
// 	book_price := ctx.Param("price")

// 	err := config.DB.QueryRow("SELECT BookId, BookName, BookPrice, PublisherId, UpdateDate FROM tb_books WHERE BookId=? AND BookPrice=?", book_id, book_price).Scan(
// 		&book.ID, &book.Name, &book.Price, &book.PublisherID, &book.UpdatedDate,
// 	)

// 	if err != nil {
// 		log.Fatalln("[SQL][Error][GetBookById] Setection Error :", err.Error())
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"book": book})
// 	log.Println("[SQL][GetBookById] SELECTED BOOK: ", book_id)
// }

// Book type
func GetBookTypes(ctx *gin.Context) {
	rows, err := config.DB.Query("SELECT BookTypeId ,BookTypeName ,BookOwnerName ,Description ,Status ,UpdateDate ,UpdateBy FROM tb_BookTypes ORDER BY BookTypeId ASC")
	if err != nil {
		log.Fatalln("[!][book/type] SELECT failed:", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var booktypes []model.BOOK_CATALOG

	for rows.Next() {
		var bt model.BOOK_CATALOG

		err := rows.Scan(&bt.ID, &bt.Name, &bt.OwenerName, &bt.Description, &bt.Status, &bt.UpdatedAt, &bt.UpdateBy)
		if err != nil {
			log.Fatalln("[!][book/types] Scan failed:", err.Error())
			return
		}
		booktypes = append(booktypes, bt)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": booktypes})
	log.Println("[SQL][book/types] SELECTED ALL BOOK TYPES")
}

func AddBookType(ctx *gin.Context) {
	prefix := ctx.Request.FormValue("prefix")
	book_type_name := ctx.Request.FormValue("typename")
	book_type_owner_name := ctx.Request.FormValue("typeowner")
	book_type_description := ctx.Request.FormValue("description")
	book_type_status := ctx.Request.FormValue("status")
	book_type_create_by := ctx.Request.FormValue("createby")
	book_type_update_by := ctx.Request.FormValue("updateby")

	// EXEC dbo.sp_BookTypeInsert 'BKT' // book grop type

	_, err := config.DB.Query("EXEC dbo.sp_BookTypeInsert ?, ?, ?, ?, ?, ?, ?",
		prefix, book_type_name, book_type_owner_name, book_type_description,
		book_type_status, book_type_create_by, book_type_update_by)

	if err != nil {
		log.Fatalln("[!] INSERT book/type failed:", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"DONE": book_type_name})
	log.Println("[+] TYPE NAME: ", book_type_name)
}

func UpdateBookType(ctx *gin.Context) {
	book_type_id := ctx.Request.FormValue("booktypeid")
	book_type_name := ctx.Request.FormValue("booktypename")
	book_type_owner_name := ctx.Request.FormValue("booktypeowner")
	book_type_description := ctx.Request.FormValue("description")
	book_type_status := ctx.Request.FormValue("status")
	book_type_update_by := ctx.Request.FormValue("updateby")

	// EXEC [dbo].[sp_BookTypeUpdate] '',''

	_, err := config.DB.Query("EXEC dbo.sp_BookTypeUpdate ?, ?, ?, ?, ?, ?", book_type_id,
		book_type_name, book_type_owner_name, book_type_description, book_type_status, book_type_update_by)

	if err != nil {
		log.Fatalln("[!] UPDATE book type failed:", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"DONE": book_type_id})
	log.Println("[*] BOOK TYPE NAME: ", book_type_id)
}

func DeleteBookType(ctx *gin.Context) {
	book_type_id := ctx.Request.FormValue("booktypeid")

	// EXEC dbo.sp_GroupInsert 'GBK','TEST', '1', 'Test how to call stored procedure', 1, 'ADMIN', 'ADMIN'

	_, err := config.DB.Query("EXEC dbo.sp_BookTypeDelete ?", book_type_id)

	if err != nil {
		log.Fatalln("[!] DELETE book type failed:", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"DONE": book_type_id})
	log.Println("[-] GROUP NAME: ", book_type_id)
	// [dev] *need update date and update by for record
}
