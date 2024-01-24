package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // nice way to import drivers
	"io"
	"strings"
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

func sampleData() []string {
	// turn this into an array
	entries := `1234567890
	123 456 7891
	(123) 456 7892
	(123) 456-7893
	123-456-7894
	123-456-7890
	1234567892
	(123)456-7892`

	result := make([]string, 0)
	reader := bufio.NewReader(strings.NewReader(entries))
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		line = strings.TrimSpace(line)
		result = append(result, line)
		//fmt.Println(line)
	}
	return result
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
	for _, number := range sampleData() {
		_, err = db.Exec("INSERT INTO phone_numbers (NUMBER) values ($1)", number)
		if err != nil {
			panic(err)
		}
	}
}
