package main

import (
	"fmt"
)

func main() {
	var str = "abcde😎fg"
	fmt.Printf("String = %v (size %d)\n", str, len(str))
	fmt.Println("Loop bytes")
	for i, c := range []byte(str) {
		fmt.Printf("[%d]=%v\n", i, c)
	}
	fmt.Println("Loop runes")
	for i, c := range str {
		fmt.Printf("[%d]=%v\n", i, c)
	}
	/*
		String = abcde😎fg (size 11)

		Loop bytes
		[0]=97
		[1]=98
		[2]=99
		[3]=100
		[4]=101
		[5]=240
		[6]=159
		[7]=152
		[8]=142
		[9]=102
		[10]=103

		Loop runes
		[0]=97
		[1]=98
		[2]=99
		[3]=100
		[4]=101
		[5]=128526 << this emoticon takes 4 bytes!
		[9]=102
		[10]=103
	*/
}
