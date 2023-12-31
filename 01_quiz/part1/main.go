package main

import (	
	"fmt"
    "flag"
	"bufio"
	"log"
	"os"
	"strings"
	"encoding/csv"
) 
func main () {
	csvFile := flag.String("csvFile", "records.csv", "CSV file containing quiz items")
	flag.Parse()

	data := ReadCSV(*csvFile)
	rights, wrongs, totals := runQuiz(data)
    
	fmt.Printf("You got %d out of %d! (%d wrong) \n", rights, totals, wrongs)
}

func ReadCSV(fileName string) (data [][]string) {
	f,err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err = csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func runQuiz(items [][]string) (rights, wrongs, totals int) {
	reader := bufio.NewReader(os.Stdin)
	for _, item := range items {
       question := item[0]
	   answer   := item[1]
	   fmt.Printf("%v ? ", question)
	   got, _ := reader.ReadString('\n')
	   got = strings.Replace(got, "\n", "", -1)

	   if strings.Compare(got, answer) == 0{
		 rights++
	   } else {
		 wrongs++
	   }	
	}

	return rights, wrongs, len(items)
}