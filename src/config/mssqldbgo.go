package config

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

var Sever = "159.65.5.144"
var Port = 1433
var User = "sa"
var Password = "Penbun@1234"
var Database = "PENBUN"

// DB database global
var DB *sql.DB
var ConnectionString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", Sever, User, Password, Port, Database)
