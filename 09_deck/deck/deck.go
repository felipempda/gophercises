package deck

//go:generate go run golang.org/x/tools/cmd/stringer -type=Suit,Rank

import (
	"fmt"
	"math/rand"
	"sort"
)

type Suit uint8
type Rank uint8 // smaller than int

const (
	Spade   Suit = iota // "♠"
	Diamond             //"♦"
	Club                //"♣"
	Heart               //"♥"
	Joker
)

const (
	_ Rank = iota // the first is not used, so actually we start from 1
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

type Card struct {
	Suit
	Rank
}

func translateSuit(s Suit) string {
	switch s {
	case Spade:
		return "♠"
	case Diamond:
		return "♦"
	case Club:
		return "♣"
	case Heart:
		return "♥"
	default:
		return ""
	}
}

// It can't be a * argument for the String()!
func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)
}

type Deck []Card

func New(opts ...func(Deck) Deck) Deck {
	var result Deck
	for suit := Spade; suit <= Heart; suit = suit + 1 {
		for face := Ace; face <= King; face = face + 1 {
			result = append(result, Card{suit, face})
		}
	}

	// run the function for each option
	for _, opt := range opts {
		result = opt(result)
	}
	return result
}

func (d Deck) Len() int {
	return len(d)
}

func (d Deck) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// default Sort: by Rank and then by Suit
func (d Deck) Less(i, j int) bool {
	if d[i].Rank != d[j].Rank {
		return d[i].Rank < d[j].Rank
	}
	if d[i].Suit != d[j].Suit {
		return d[i].Suit < d[j].Suit
	}
	return false
}

// CustomSort:

// By is the type of a "less" function that defines the ordering of its Cards arguments.
type By func(c1, c2 *Card) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(cards []Card) {

	cs := &cardSorter{
		cards: cards,
		by:    by,
	}

	sort.Sort(cs)
}

// cardSorter joins a By function and a slice of Planets to be sorted.
type cardSorter struct {
	cards []Card
	by    By
}

// Len is part of sort.Interface.
func (s *cardSorter) Len() int {
	return len(s.cards)
}

// Swap is part of sort.Interface.
func (s *cardSorter) Swap(i, j int) {
	s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *cardSorter) Less(i, j int) bool {
	return s.by(&s.cards[i], &s.cards[j])
}

// another way of doing it, with functional options pattern
// (it looks like Java BuilderPattern, but less organized - what are actually the options?)

// absolute Value of Value of Card , convert to single int the ideal order
func AbsRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

// this is one of the "options", it implements the opts signature
func DefaultSort(deck Deck) Deck {
	sort.Slice(deck, Less(deck))
	return deck
}

// this implements Less interface, required by sort.Slice
func Less(deck Deck) func(i, j int) bool {
	return func(i, j int) bool {
		return AbsRank(deck[i]) < AbsRank(deck[j])
	}
}

// now for a generic order solution with functional arguments

// Sort receives a function that returns a function
// and returns a function
func Sort(less func(Deck) func(i, j int) bool) func(Deck) Deck {
	return func(deck Deck) Deck {
		sort.Slice(deck, less(deck))
		return deck
	}
}

func Shuffle(deck Deck) Deck {
	swap := func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	}
	rand.Shuffle(len(deck), swap)
	return deck
}
