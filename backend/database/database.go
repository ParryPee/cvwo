package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(username string, password string, host string, port int, dbName string) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", username, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error validating sql.Open arguments: ", err)
	}
	return db
}
