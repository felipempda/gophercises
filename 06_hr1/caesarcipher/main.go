package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter length of string: ")
	length, _ := strconv.Atoi(readStr(reader))
	fmt.Print("Enter string: ")
	Str := readStr(reader)
	fmt.Print("Enter number of entries to shift: ")
	shift, _ := strconv.Atoi(readStr(reader))
	newString := CaesarCipher(Str, length, shift)
	fmt.Println(newString)

}

func readStr(reader *bufio.Reader) string {
	got, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	got = strings.Replace(got, "\n", "", -1)
	return got
}

func CaesarCipher(Str string, length, shift int) string {
	fmt.Println("String: ", Str)
	fmt.Println("Shift: ", shift)
	newString := make([]byte, len(Str))
	var newAscii int
	for i, c := range Str {
		ascii := int(c)
		if ascii >= 65 && ascii <= 90 { //ascii 65-90:= A-Z
			newAscii = ascii + (shift % 26)
			if newAscii > 90 {
				newAscii = newAscii - 26
			}
		} else if ascii >= 97 && ascii <= 122 { //ascii 97-122:= a-z
			newAscii = ascii + (shift % 26)
			if newAscii > 122 {
				newAscii = newAscii - 26
			}
		} else {
			newAscii = ascii
		}
		fmt.Printf("Current Ascii Value: %v, New: %v\n", ascii, newAscii)
		newString[i] = byte(newAscii)
	}
	return string(newString)
}
