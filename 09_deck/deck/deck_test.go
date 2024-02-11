package deck

import (
	"fmt"
	"sort"
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

func ExampleSort() {
	cards := Deck{
		Card{Heart, Ten},
		Card{Heart, Nine},
		Card{Spade, Ace},
		Card{Spade, Two},
	}

	sort.Sort(cards)
	for _, card := range cards {
		fmt.Println(card)
	}

	// Output:
	// Ace of Spades
	// Two of Spades
	// Nine of Hearts
	// Ten of Hearts
}

func ExampleCustomSort() {
	cards := Deck{
		Card{Heart, Ten},
		Card{Heart, Nine},
		Card{Spade, Ace},
		Card{Spade, Two},
	}

	rankDescended := func(c1, c2 *Card) bool {
		return c1.Rank > c2.Rank
	}

	By(rankDescended).Sort(cards)
	for _, card := range cards {
		fmt.Println(card)
	}

	// Output:
	// Ten of Hearts
	// Nine of Hearts
	// Two of Spades
	// Ace of Spades
}
