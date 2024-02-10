package deck

//go:generate go run golang.org/x/tools/cmd/stringer -type=Suit,Rank

import (
	"fmt"
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
	One
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

func New() []Card {
	result := make([]Card, 0)
	for suit := Spade; suit <= Heart; suit = suit + 1 {
		for face := Ace; face <= King; face = face + 1 {
			c := Card{suit, face}
			result = append(result, c)
		}
	}
	return result
}
