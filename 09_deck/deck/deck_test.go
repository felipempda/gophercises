package deck

import (
	"fmt"
	"testing"
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

func TestNew(t *testing.T) {
	cards := New()

	l := len(cards)
	if l != 52 {
		t.Errorf("Wrong number of cards in the deck: %d", l)
	}
}
