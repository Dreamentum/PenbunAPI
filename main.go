package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

var server = "139.59.228.81"
var port = 1433
var user = "sa"
var password = "admin@1q2w3e4r5t"
var database = "PNB"

var db *sql.DB
var err error

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)

	db, err = sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("[MSSQL][Error] Open connection failed:", err.Error())
	}
	log.Printf("[MSSQL] Connected!\n")
	defer db.Close()

	checkMssqlVersion()

	checkDbVersion()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/books", getAllBooks).Methods("GET")
	router.HandleFunc("/api/v1/books/{id}", getBookById).Methods("GET")
	router.HandleFunc("/api/v1/books", createBook).Methods("POST")
	router.HandleFunc("/api/v1/books/{id}", updateBookById).Methods("PUT")
	router.HandleFunc("/api/v1/books/{id}", deleteBookById).Methods("DELETE")

	log.Printf("[GOLANG] Running PENBUN API at api.penbun.com:8080\n")
	http.ListenAndServe(":8080", router)
}

func checkMssqlVersion() {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("[MSSQL][Error][checkMssqlVersion] Ping database server failed:", err.Error())
	}

	var version string

	err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&version)
	if err != nil {
		log.Fatal("[MSSQL][Error][checkMssqlVersion] Scan failed:", err.Error())
	}
	log.Printf("[MSSQL] version %s\n", version)
}

func checkDbVersion() {
	ctx := context.Background()
	var version string

	err = db.QueryRowContext(ctx, "SELECT TOP(1) version_no FROM dbo.version ORDER BY id DESC").Scan(&version)
	if err != nil {
		log.Fatal("[MSSQL][Error][checkDbVersion] Scan failed:", err.Error())
	}
	log.Printf("[MSSQL] PEBUN DATABASE %s\n", version)
}

type book struct {
	Id         string  `json:"id"`
	Isbn       string  `json:"isbn"`
	Name       string  `json:"name"`
	AuthorName string  `json:"authorName"`
	Year       string  `json:"year"`
	Price      float64 `json:"price"`
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT BookId, BookName, PublisherId, UpdateDate, BookPrice FROM tb_books ORDER BY BookId ASC")
	if err != nil {
		log.Fatal("[MSSQL][Error][getAllBooks] SELECT failed:", err.Error())
	}
	defer rows.Close()

	var books []book
	for rows.Next() {
		var b book
		err := rows.Scan(&b.Id, &b.Name, &b.AuthorName, &b.Year, &b.Price)
		if err != nil {
			log.Fatal("[MSSQL][Error][getAllBooks] Scan failed:", err.Error())
		}
		books = append(books, b)
	}
	json.NewEncoder(w).Encode(books)
	log.Printf("[MSSQL][getAllBooks] SELECTED ALL BOOKS\n")
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	rows, err := db.Query("SELECT BookId, BookName, PublisherId, UpdateDate, BookPrice FROM tb_books WHERE BookId=?", param["id"])
	if err != nil {
		log.Fatal("[MSSQL][Error][getBookById] SELECT failed:", err.Error())
	}
	defer rows.Close()

	var b book
	for rows.Next() {
		err := rows.Scan(&b.Id, &b.Name, &b.AuthorName, &b.Year, &b.Price)
		if err != nil {
			log.Fatal("[MSSQL][Error][getBookById] Scan failed:", err.Error())
		}
	}
	json.NewEncoder(w).Encode(b)
	log.Printf("[MSSQL][Error][getBookById] SELECTED A BOOK %s\n", b.Name)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var b book
	json.NewDecoder(r.Body).Decode(&b)

	stmt, err := db.Prepare("INSERT INTO book(isbn, name, authorName, year, totalPage) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("[MSSQL][Error][createBook] Prepare failed:", err.Error())
	}

	_, err = stmt.Exec(b.Isbn, b.Name, b.AuthorName, b.Year, b.Price)
	if err != nil {
		log.Fatal("[MSSQL][Error][createBook] INSERT failed:", err.Error())
	}
	defer stmt.Close()
	w.WriteHeader(http.StatusCreated)
	log.Printf("[MSSQL][createBook] INSERTED A NEW BOOK %s\n", b.Name)
}

func updateBookById(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	var b book
	json.NewDecoder(r.Body).Decode(&b)

	stmt, err := db.Prepare("UPDATE book SET isbn=?, name=?, authorName=?, year=?, totalPage=? WHERE id=?")
	if err != nil {
		log.Fatal("[MSSQL][Error][updateBookById] Prepare failed:", err.Error())
	}

	_, err = stmt.Exec(b.Isbn, b.Name, b.AuthorName, b.Year, b.Price, param["id"])
	if err != nil {
		log.Fatal("[MSSQL][Error][updateBookById] UPDATE failed:", err.Error())
	}
	defer stmt.Close()
	w.WriteHeader(http.StatusOK)
	log.Printf("[MSSQL][updateBookById] UPDATED A BOOK %s\n", b.Name)
}

func deleteBookById(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM book WHERE id=?")
	if err != nil {
		log.Fatal("[MSSQL][deleteBookById] Prepare failed:", err.Error())
	}

	_, err = stmt.Exec(param["id"])
	if err != nil {
		log.Fatal("[MSSQL][deleteBookById] Database DELETE failed:", err.Error())
	}
	defer stmt.Close()
	w.WriteHeader(http.StatusOK)
	log.Printf("[MSSQL][deleteBookById] DELETED A BOOK \n")
}
