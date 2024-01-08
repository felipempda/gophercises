package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
 * Complete the 'caesarCipher' function below.
 *
 * The function is expected to return a STRING.
 * The function accepts following parameters:
 *  1. STRING s
 *  2. INTEGER k
 */

func caesarCipher(s string, i, k int32) string {
	// Write your code here
	shift := int(k)
	newString := make([]byte, len(s))
	var newAscii int
	var newChar byte
	for i, c := range s {
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
		//fmt.Printf("Current Ascii Value: %v, New: %v\n", ascii, newAscii)
		newChar = byte(newAscii)
		newString[i] = newChar
	}
	return string(newString)
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	nTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	n := int32(nTemp)

	s := readLine(reader)

	kTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	k := int32(kTemp)

	result := caesarCipher(s, n, k)

	fmt.Fprintf(writer, "%s\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
