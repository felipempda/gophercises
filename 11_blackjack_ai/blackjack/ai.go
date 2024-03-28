package blackjack

import (
	"bufio"
	"fmt"
	"github.com/felipempda/gophercises/09_deck/deck"
	"os"
	"strings"
)

// interface
type AI interface {
	Bet(shuffled bool) int
	Play(handIndx int, hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card)
}

// human implementation of interface
func HumanAI() AI {
	return humanAI{}
}

type humanAI struct {
}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("Deck has just been shuffled!")
	}
	fmt.Printf("How much would you like to bet? ")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

func (ai humanAI) Play(handIndx int, hand []deck.Card, dealer deck.Card) Move {
	fmt.Printf("[ YOUR TURN - HAND %d ]\n", handIndx)
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf("Dealer: %s = %d \n", dealer, Score(dealer))
	fmt.Printf("Player [hand %d] (%d cards): %s = %d \n", handIndx, len(hand), handDeck(hand), Score(hand...))
	for {
		fmt.Printf("Press (h)it or (s)tand or (d)ouble or s(p)lit ")
		reader := bufio.NewReader(os.Stdin)
		got, _ := reader.ReadString('\n')
		got = strings.Replace(got, "\n", "", -1)

		switch got {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		case "p":
			return MoveSplit
		default:
			fmt.Println("Wrong answer, try again!")
		}
	}
}

func (ai humanAI) Results(player [][]deck.Card, dealer []deck.Card) {
	fmt.Printf("[ HAND  - RESULTS ]\n")
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf(" Dealer (%d cards): %s = %d \n", len(dealer), handDeck(dealer), Score(dealer...))
	for idx, hand := range player {
		fmt.Printf(" Player [hand %d] (%d cards): %s = %d \n", idx+1, len(hand), handDeck(hand), Score(hand...))
	}
	fmt.Printf("------------------------------------------------------------------------------\n")
}
