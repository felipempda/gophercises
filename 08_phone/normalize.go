package main

import (
	"fmt"
	"github.com/felipempda/gophercises/08_phone/phone"
)

func main() {
	for _, number := range phone.SampleData() {
		fmt.Printf("Number=%s, Normalized=%s\n", number, phone.Normalize(number))

	}
}
