package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(username string, password string) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/cvwo_forum?parseTime=true", username, password)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error validating sql.Open arguments: ", err)
	}
	return db
}
