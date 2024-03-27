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
	Bet() int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card)
}

// human implementation of interface
func HumanAI() AI {
	return humanAI{}
}

type humanAI struct {
}

func (ai humanAI) Bet() int {
	return 1
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	fmt.Printf("Dealer: %s = %d \n", dealer, Score(dealer))
	fmt.Printf("Player (%d cards): %s = %d \n\n", len(hand), handDeck(hand), Score(hand...))
	var decided Move
	for true {
		fmt.Printf("Press (h)it or (s)tand ")
		reader := bufio.NewReader(os.Stdin)
		got, _ := reader.ReadString('\n')
		got = strings.Replace(got, "\n", "", -1)

		if got == "h" {
			decided = MoveHit
			break
		} else if got == "s" {
			decided = MoveStand
			break
		} else {
			fmt.Println("Wrong answer, try again!")
		}
	}
	return decided
}

func (ai humanAI) Results(player [][]deck.Card, dealer []deck.Card) {
	dealer_points := Score(dealer...)
	player_points := Score(player[0]...)

	fmt.Printf("[ BLACK JACK  - RESULTS ]\n\n")
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf(" Dealer (%d cards): %s = %d \n", len(dealer), handDeck(dealer), dealer_points)
	fmt.Printf(" Player (%d cards): %s = %d \n ", len(player[0]), handDeck(player[0]), player_points)
	fmt.Printf("------------------------------------------------------------------------------\n")
}
