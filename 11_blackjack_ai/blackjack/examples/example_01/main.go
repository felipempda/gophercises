package main

import (
	"fmt"
	"github.com/felipempda/gophercises/09_deck/deck"
	"github.com/felipempda/gophercises/11_blackjack_ai/blackjack"
)

func main() {
	ai := basicAI{
		decks: 4,
	} // blackjack.HumanAI()
	opts := blackjack.Options{
		Hands: 5000,
		Decks: 4,
	}
	game := blackjack.New(opts)
	winnings := game.PlayGame(&ai)
	fmt.Println("Our AI won/lost:", winnings)
}

type basicAI struct {
	score int
	seen  int
	decks int
}

func (ai *basicAI) Bet(shuffled bool) int {
	if shuffled {
		ai.score = 0
		ai.seen = 0
	}
	trueScore := ai.score / ((ai.decks*52 - ai.seen) / 52)
	if trueScore > 10 {
		fmt.Println(ai)
	}
	switch {
	case trueScore > 14:
		return 10000
	case trueScore > 8:
		return 500
	default:
		return 100
	}
}

func (ai *basicAI) Play(handIndx int, hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	if score == 21 {
		return blackjack.MoveStand
	}
	if len(hand) == 2 {
		if hand[0] == hand[1] {
			cardScore := blackjack.Score(hand[0])
			if cardScore >= 8 && cardScore != 10 {
				return blackjack.MoveSplit
			}
		}
		if (score == 10 || score == 11) && !blackjack.Soft(hand...) {
			return blackjack.MoveDouble
		}
	}
	dealerScore := blackjack.Score(dealer)
	if dealerScore >= 5 && dealerScore <= 6 {
		return blackjack.MoveStand
	}
	if score < 13 {
		return blackjack.MoveHit
	}
	return blackjack.MoveStand
}

func (ai *basicAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	fmt.Printf("[ HAND  - RESULTS ]\n")
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf(" Dealer (%d cards): %s = %d \n", len(dealer), (dealer), blackjack.Score(dealer...))
	for idx, hand := range hands {
		fmt.Printf(" Player [hand %d] (%d cards): %s = %d \n", idx+1, len(hand), (hand), blackjack.Score(hand...))
	}
	fmt.Printf("------------------------------------------------------------------------------\n")

	for _, card := range dealer {
		ai.count(card)
	}
	for _, hand := range hands {
		for _, card := range hand {
			ai.count(card)
		}
	}
}

func (ai *basicAI) count(card deck.Card) {
	score := blackjack.Score(card)
	switch {
	case score >= 10:
		ai.score--
	case score <= 6:
		ai.score++
	}
	ai.seen++
}
