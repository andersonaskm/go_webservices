package database

import (
	"database/sql"
	"log"
	"time"
)

var DbConn *sql.DB

func SetUpDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", "inventorydb:M3s9H67F-Cy-/$N@tcp(127.0.0.1:3306)/inventorydb")
	if err != nil {
		log.Fatal(err)
	}
	DbConn.SetMaxOpenConns(10)
	DbConn.SetMaxIdleConns(10)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}
