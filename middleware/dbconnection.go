package middleware

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// create connection with postgres db
func CreateConnection() *sql.DB {
	//Load .env variable
	conn := DotEnvVariable("MYSQL_URL")
	// Open the connection
	db, err := sql.Open("mysql", conn)

	if err != nil {
		panic(err)
	} else {
		log.Println("connected to MySql Database..")
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	// return the connection
	return db
}
