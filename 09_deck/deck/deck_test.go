package deck

import (
	"fmt"
)

func ExampleCards() {

	fmt.Println(Card{Heart, Ten})
	fmt.Println(Card{Spade, Ace})
	fmt.Println(Card{Club, Queen})
	fmt.Println(Card{Diamond, Two})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ten of Hearts
	// Ace of Spades
	// Queen of Clubs
	// Two of Diamonds
	// Joker
}
