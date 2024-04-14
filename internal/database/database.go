package database

import (
	"database/sql"
	"log"
)

func InitDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/test-database?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	{
		query := `
			CREATE TABLE IF NOT EXISTS users (
				id INT AUTO_INCREMENT,
				username TEXT NOT NULL,
				password TEXT NOT NULL,
				created_at DATETIME,
				PRIMARY KEY (id)
			);`
	
		if _, err := db.Exec(query); err != nil {
				log.Fatal(err)
		}
	}

	return db
}