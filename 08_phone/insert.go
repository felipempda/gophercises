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
	//dbName   = "postgres"
)

func main() {

	populateDatabase()

}

func populateDatabase() {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("drop table phone_numbers") // as a former DBA this feels so right 😈
	_, err = db.Exec("CREATE TABLE phone_numbers(ID SERIAL PRIMARY KEY, NUMBER TEXT NOT NULL)")
	for _, number := range phone.SampleData() {
		_, err = db.Exec("INSERT INTO phone_numbers (NUMBER) values ($1)", number)
		if err != nil {
			panic(err)
		}
	}
}
