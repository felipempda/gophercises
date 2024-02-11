package main

import (
	"fmt"
	"github.com/felipempda/gophercises/09_deck/deck"
	"sort"
)

func main() {

	//test simple
	fmt.Println(deck.Spade)
	c := deck.Card{deck.Spade, deck.Ten}
	fmt.Println(c)
	fmt.Println(c)

	c.Suit = c.Suit + 1
	fmt.Println(c)

	c.Suit = c.Suit + 1
	fmt.Println(c)

	c.Suit = c.Suit + 1
	fmt.Println(c)

	for suit := deck.Spade; suit <= deck.Heart; suit = suit + 1 {
		fmt.Println(suit)

	}

	// test New Deck
	myDeck := deck.New()
	for _, card := range myDeck {
		fmt.Println(card)
	}

	// testingSort
	fmt.Printf("Size of deck: %d\n", len(myDeck))

	sort.Sort(myDeck)
	fmt.Println("Sorted deck...")
	fmt.Println(myDeck)
}
