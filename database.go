package main

import (

	// Adjust the import path as necessary
	"Student_RESTAPI/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	// dsn := "root:@Surya123@tcp(localhost:3306)/rest_api?parseTime=true"
	// log.Printf("Connecting to database with DSN: %s", dsn)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	log.Println("Database connection established successfully.")
}
