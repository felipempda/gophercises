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

	// return this order of Suits Heart, Spade, Club, Diamond, descended Ranks
	customRank := func(c deck.Card) int {
		r := 0
		switch c.Suit {
		case deck.Heart:
			r = 100
		case deck.Spade:
			r = 200
		case deck.Club:
			r = 300
		case deck.Diamond:
			r = 400
		default:
			r = 500
		}
		return r - int(c.Rank)
	}
	ByCustomRank := func(d deck.Deck) func(i, j int) bool {
		return func(i, j int) bool {
			return customRank(d[i]) < customRank(d[j])
		}
	}
	fmt.Println("Custom Sorted deck...")
	newDeck := deck.New(deck.Sort(ByCustomRank))
	fmt.Println(newDeck)

	fmt.Println("Shuffled deck...")
	shuffledDeck := deck.New(deck.Shuffle)
	fmt.Println(shuffledDeck)

	fmt.Println("With Jokers deck...")
	withJokers := deck.New(deck.WithJokers(3))
	fmt.Println(withJokers)

	fmt.Println("Remove Jacks, Queens and Kings...")
	removeJQK := func(c deck.Card) bool {
		return c.Rank == deck.Jack || c.Rank == deck.Queen || c.Rank == deck.King
	}
	filteredDeck := deck.New(deck.Filter(removeJQK))
	fmt.Println(filteredDeck)

	fmt.Println("TwoDecks...")
	twoDeck := deck.New(deck.WithMultipleDecks(3))
	fmt.Println(twoDeck)
	fmt.Println("Size: ", len(twoDeck))

	fmt.Println("TwoDecks + Filtered...")
	twoDeckFilter := deck.New(deck.WithMultipleDecks(3), deck.Filter(removeJQK))
	fmt.Println(twoDeckFilter)
	fmt.Println("Size: ", len(twoDeckFilter))
}
