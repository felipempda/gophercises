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
	fmt.Println("Inserting...")
	populateDatabase(db)

	fmt.Println("Selecting...")
	must(queryData(db))

	fmt.Println("Updating...")
	updateData(db)

	fmt.Println("Selecting...")
	must(queryData(db))

	fmt.Println("Deleting duplicates...")
	deleteDuplicate(db)

	fmt.Println("Selecting...")
	must(queryData(db))
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

type phoneData struct {
	id     int
	number string
}

func selectDataGeneral(db *sql.DB, orderBy string) ([]phoneData, error) {
	rows, err := db.Query("SELECT ID, NUMBER FROM PHONE_NUMBERS ORDER BY " + orderBy)
	var results []phoneData
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var p phoneData
		err = rows.Scan(&p.id, &p.number)
		if err != nil {
			return results, err
		}
		results = append(results, p)
		//fmt.Printf("%d - %s\n", p.id, p.number)
	}

	if err = rows.Err(); err != nil {
		return results, err
	}

	return results, err
}

func selectDataOrderedByID(db *sql.DB) ([]phoneData, error) {
	return selectDataGeneral(db, "ID")
}

func selectDataOrderedByNumber(db *sql.DB) ([]phoneData, error) {
	return selectDataGeneral(db, "NUMBER")
}

func queryData(db *sql.DB) error {
	phones, err := selectDataOrderedByID(db)
	if err != nil {
		return err
	}
	for _, p := range phones {
		fmt.Printf("%d - %s\n", p.id, p.number)
	}
	return nil
}

func updateData(db *sql.DB) {
	phones, err := selectDataOrderedByID(db)
	must(err)
	for _, p := range phones {
		newNumber := phone.Normalize(p.number)
		if p.number != newNumber {
			_, err := db.Exec("UPDATE PHONE_NUMBERS SET NUMBER = $1 WHERE ID = $2", newNumber, p.id)
			must(err)
		}
	}
}

func deleteDuplicate(db *sql.DB) {
	phones, err := selectDataOrderedByNumber(db)
	must(err)
	var previousNumber string
	for k, p := range phones {
		if previousNumber == p.number && k > 0 {
			_, err := db.Exec("DELETE FROM PHONE_NUMBERS WHERE ID  = $1", p.id)
			must(err)
		}
		previousNumber = p.number
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
