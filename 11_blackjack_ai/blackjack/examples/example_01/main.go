package main

import (
	//"github.com/felipempda/gophercises/09_deck/deck"
	"fmt"
	"github.com/felipempda/gophercises/11_blackjack_ai/blackjack"
)

func main() {
	ai := blackjack.HumanAI()
	opts := blackjack.Options{
		Hands: 100,
		Decks: 3,
	}
	game := blackjack.New(opts)
	winnings := game.PlayGame(ai)
	fmt.Println("Our AI won/lost:", winnings)
}
