package main

import (
	"database/sql"
	"fmt"
	"github.com/felipempda/gophercises/08_phone/phone"
	_ "github.com/lib/pq" // nice way to import drivers
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass123"
	dbName   = "phone_database"
)

func main() {

	db, err := createDB()
	must(err)
	defer db.Close()
	populateDatabase(db)

}

func createDB() (*sql.DB, error) {
	// connect first time
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", dbInfo)
	must(err)

	// recreateDB
	err = recreateDB(db, dbName)
	must(err)
	db.Close()

	// connect to database
	dbInfo = fmt.Sprintf("%s dbname=%s", dbInfo, dbName)
	db, err = sql.Open("postgres", dbInfo)
	db.Ping()

	// return connection
	return db, err
}

func populateDatabase(db *sql.DB) {
	_, err := db.Exec("drop table phone_numbers") // as a former DBA this feels so right 😈
	_, err = db.Exec("CREATE TABLE phone_numbers(ID SERIAL PRIMARY KEY, NUMBER TEXT NOT NULL)")
	must(err)
	for _, number := range phone.SampleData() {
		// result, err := db.Exec("INSERT INTO phone_numbers (NUMBER) values ($1)", number)
		// id, err := result.LastInsertId()       // postgres doesn't support: panic: LastInsertId is not supported by this driver
		var id int
		err := db.QueryRow("INSERT INTO phone_numbers (NUMBER) values ($1) RETURNING id", number).Scan(&id)
		fmt.Printf("id=%d\n", id)
		must(err)
	}
}

func recreateDB(db *sql.DB, dbName string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE DATABASE " + dbName)
	return err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
