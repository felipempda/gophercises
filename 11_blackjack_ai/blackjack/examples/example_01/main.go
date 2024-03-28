package main

import (
	"fmt"
	"github.com/felipempda/gophercises/09_deck/deck"
	"github.com/felipempda/gophercises/11_blackjack_ai/blackjack"
)

func main() {
	ai := blackjack.HumanAI()
	opts := blackjack.Options{
		Hands: 3,
		Decks: 3,
	}
	game := blackjack.New(opts)
	winnings := game.PlayGame(ai)
	fmt.Println("Our AI won/lost:", winnings)
}

type basicAI struct {
}

func (ai *basicAI) Bet(shuffled bool) int {
	panic("not implemented") // TODO: Implement
}

func (ai *basicAI) Play(handIndx int, hand []deck.Card, dealer deck.Card) blackjack.Move {
	panic("not implemented") // TODO: Implement
}

func (ai *basicAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	panic("not implemented") // TODO: Implement
}
