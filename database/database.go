package database

import (
	"database/sql"
)

var Db *sql.DB

func Connect(db *sql.DB) {
	Db = db
}

func Query(q string) int64 {
	stmtOut, err := Db.Prepare(q)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var squareNum int64 // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return squareNum
}
