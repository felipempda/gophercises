package deck

import (
	"fmt"
)

type Suit int
type Face int

const (
	Spade   Suit = iota // "♠"
	Diamond             //"♦"
	Club                //"♣"
	Heart               //"♥"
)

const (
	CardA Face = iota
	Card1
	Card2
	Card3
	Card4
	Card5
	Card6
	Card7
	Card8
	Card9
	Card10
	CardJ
	CardQ
	CardK
)

type Card struct {
	Suit Suit
	Face Face
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

func translateFace(c Face) string {
	switch c {
	case CardA:
		return "A"
	case CardJ:
		return "J"
	case CardQ:
		return "Q"
	case CardK:
		return "K"
	default:
		return fmt.Sprintf("%d", c)
	}
}

func (c *Card) Describe() string {
	return fmt.Sprintf("%s%s", translateFace(c.Face), translateSuit(c.Suit))
}

func New() []Card {
	result := make([]Card, 0)
	for suit := Spade; suit <= Heart; suit = suit + 1 {
		for face := CardA; face <= CardK; face = face + 1 {
			c := Card{suit, face}
			result = append(result, c)
		}
	}
	return result
}
