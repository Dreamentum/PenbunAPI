package config

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

var Sever = "134.209.104.247"
var Port = 1433
var User = "sa"
var Password = "penbun@1q2w3e4r5t"
var Database = "PNB"

// DB database global
var DB *sql.DB
var ConnectionString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", Sever, User, Password, Port, Database)
