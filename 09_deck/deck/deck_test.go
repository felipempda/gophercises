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

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)

	if e, g := (Card{Rank: Ace, Suit: Spade}), (cards[0]); e != g {
		t.Errorf("Wrong first card, expected %s got %s", e, g)
	}
	if e, g := (Card{Rank: King, Suit: Heart}), (cards[len(cards)-1]); e != g {
		t.Errorf("Wrong last card, expected %s got %s", e, g)
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))

	if e, g := (Card{Rank: Ace, Suit: Spade}), (cards[0]); e != g {
		t.Errorf("Wrong first card, expected %s got %s", e, g)
	}
	if e, g := (Card{Rank: King, Suit: Heart}), (cards[len(cards)-1]); e != g {
		t.Errorf("Wrong last card, expected %s got %s", e, g)
	}
}

func TestWithJokers(t *testing.T) {
	jokers := 2
	cards := New(WithJokers(jokers))

	found := 0
	for _, c := range cards {
		if c.Suit == Joker {
			found++
		}
	}
	if found != jokers {
		t.Errorf("Wrong number of jokers, expected %d, found %d", jokers, found)
	}
}

func TestFilter(t *testing.T) {
	removeJQK := func(c Card) bool {
		return c.Rank == Jack || c.Rank == Queen || c.Rank == King
	}

	filteredDeck := New(Filter(removeJQK))

	for i, c := range filteredDeck {
		if removeJQK(c) == true {
			t.Errorf("Wrong filter function, filtered card still on the deck[%d]: %s", i, c)
		}
	}

}

func TestWithMultipleDecks(t *testing.T) {
	single := New()
	x := 3
	multiple := New(WithMultipleDecks(x))

	if e, f := len(single)*x, len(multiple); e != f {
		t.Errorf("Wrong number of cards with multipleDecks option. Expected %d, Found %d cards", e, f)
	}
}
