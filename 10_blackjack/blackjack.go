package main

import (
	"bufio"
	"fmt"
	"github.com/felipempda/gophercises/09_deck/deck"
	"os"
	"os/exec"
	"strings"
)

type decision int

const (
	stand decision = iota
	hit
)

type participant struct {
	cards  []deck.Card
	points int
	stand  bool
}

var (
	gameRound        int
	gameDeck         deck.Deck
	gameDeckPosition int
	player1          participant
	dealer           participant
)

func main() {

	play()
}

func drawCard() deck.Card {
	card := gameDeck[gameDeckPosition]
	gameDeckPosition = gameDeckPosition + 1
	return card
}

func play() {
	initialDraw()
	for gameRound = 1; shouldContinue(); gameRound++ {
		printScreen(false)
		processDecision(waitPlayersDecision())
	}
	printScreen(true)
}

func initialDraw() {
	// start Deck
	gameDeck = deck.New(deck.Shuffle)

	// draw 2 cards
	for i := 0; i < 2; i++ {
		player1.cards = append(player1.cards, drawCard())
		dealer.cards = append(dealer.cards, drawCard())
	}
}

func calculatePointsAllPlayers() {
	player1.points = calculatePoints(player1, -1)
	dealer.points = calculatePoints(dealer, -1)
}

func shouldContinue() bool {
	calculatePointsAllPlayers()
	return !player1.stand && player1.points < 21 && dealer.points < 21
}

func waitPlayersDecision() decision {
	var decided decision
	for true {
		fmt.Printf("Press (1) for Hit or (2) for Stand: ")
		reader := bufio.NewReader(os.Stdin)
		got, _ := reader.ReadString('\n')
		got = strings.Replace(got, "\n", "", -1)

		if got == "1" {
			fmt.Println("You chose Hit!")
			decided = hit
			break
		} else if got == "2" {
			fmt.Println("You chose Stand!")
			decided = stand
			break
		} else {
			fmt.Println("Wrong answer, try again!")
		}
	}
	return decided
}

func calculatePoints(p participant, upToPosition int) int {
	var total int
	var numberAces int
	for position, card := range p.cards {
		// considerOnly cards up to a given position?
		if position > upToPosition && upToPosition >= 0 {
			break
		}
		var value int
		switch card.Rank {
		case deck.Ace:
			value = 1
			numberAces++
		case deck.Two:
			value = 2
		case deck.Three:
			value = 3
		case deck.Four:
			value = 4
		case deck.Five:
			value = 5
		case deck.Six:
			value = 6
		case deck.Seven:
			value = 7
		case deck.Eight:
			value = 8
		case deck.Nine:
			value = 9
		default:
			value = 10
		}
		total = total + value
	}

	// special rule for 21
	for n := 1; n <= numberAces; n++ {
		if total+(n*10) == 21 {
			total = 21
			break
		}
	}
	return total
}

func processDecision(d decision) {
	if d == hit {
		player1.cards = append(player1.cards, drawCard())
		// drawACard
		if shouldContinue() {
			dealer.cards = append(dealer.cards, drawCard())
		}
	} else {
		player1.stand = true
	}
}

func printScreen(endGame bool) {

	// clear screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// show dealer cards only at the end, at the beginning only first card is shown
	var dealer_cards []deck.Card
	var dealer_points int
	if endGame {
		dealer_cards = dealer.cards
		dealer_points = calculatePoints(dealer, -1)
	} else {
		dealer_cards = dealer.cards[0:1]
		dealer_points = calculatePoints(dealer, 0)
	}

	fmt.Printf("[ BLACK JACK   -  ROUND %d ]\n\n", gameRound)
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf(" Dealer (%d cards): %s = %d \n", len(dealer.cards), dealer_cards, dealer_points)
	fmt.Printf(" Player (%d cards): %s = %d \n ", len(player1.cards), player1.cards, player1.points)
	fmt.Printf("------------------------------------------------------------------------------\n")

	if endGame {
		var burst bool
		if player1.points > 21 {
			fmt.Println("PLAYER1 BURST")
			burst = true
		}
		if dealer.points > 21 {
			fmt.Println("DEALER BURST")
			burst = true
		}
		if player1.points == 21 || (!burst && player1.points > dealer.points) {
			fmt.Println("PLAYER1 WINS!")
		} else if dealer.points == 21 || (!burst && dealer.points > player1.points) {
			fmt.Println("DEALER WINS!")
		}
	}
}
