package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Enter camelCase string:")
	strInput, _ := reader.ReadString('\n')
	fmt.Println(camelCaseCount(strInput))
}

func camelCaseCount(strInput string) int {
	words := 1
	for _, c := range strInput {
		if c != unicode.ToLower(c) {
			words = words + 1
		}
	}
	return words
}
