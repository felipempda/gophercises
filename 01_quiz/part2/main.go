package main

import (	
	"fmt"
    "flag"
	"bufio"
	"log"
	"os"
	"strings"
	"encoding/csv"
	"time"
) 
func main () {
	csvFile := flag.String("csvFile", "records.csv", "CSV file containing quiz items")
	timeout := flag.Int("timeout", 10, "Seconds to timeout quiz")
	flag.Parse()

	data := ReadCSV(*csvFile)

	rights, wrongs, totals := runQuiz(data, *timeout)
	fmt.Printf("\nYou got %d out of %d! (%d wrong) \n", rights, totals, wrongs)
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

func getInput(input chan string) {
    for {
        reader := bufio.NewReader(os.Stdin)
		got, _ := reader.ReadString('\n')
		got = strings.Replace(got, "\n", "", -1)
        input <- got
    }
}

func runQuiz(items [][]string, timeout int) (rights, wrongs, totals int) {
	input := make(chan string, 1)
	timer1 := time.NewTicker(time.Second * time.Duration(timeout))

	L:
	for _, item := range items {
       question := item[0]
	   answer   := item[1]
	   fmt.Printf("%v ? ", question)
	   go getInput(input)

	   select {

	    case got := <- input:
			if strings.Compare(got, answer) == 0{
				rights++
			} else {
				wrongs++
			}	
		case <-timer1.C:
			fmt.Printf("\n[Timeout, %d seconds elapsed!]\n", timeout)
			break L
	   }
	}

	return rights, wrongs, len(items)
}