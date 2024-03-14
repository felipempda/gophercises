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

type Hand []deck.Card

type participant struct {
	cards  Hand
	points int
	stand  bool
}

var (
	gameRound int
	gameDeck  deck.Deck
	player1   participant
	dealer    participant
)

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range strs {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func main() {

	play()
}

func drawCard() deck.Card {
	var card deck.Card
	card, gameDeck = gameDeck[0], gameDeck[1:]
	return card
}

func play() {
	//for i := 1; i<=4; i++ {
	initialDraw()
	for gameRound = 1; shouldContinue(); gameRound++ {
		printScreen(false)
		if !player1.stand {
			playersTurn()
		} else {
			dealersTurn()
		}
	}
	printScreen(true)
	// }
}

func initialDraw() {
	// start Deck
	gameDeck = deck.New(deck.WithMultipleDecks(2), deck.Shuffle)
	player1.cards = nil
	player1.stand = false
	dealer.cards = nil
	dealer.stand = false

	// draw 2 cards
	for i := 0; i < 2; i++ {

		// array of Hand, so we don't repeat code
		for _, hand := range []*Hand{&player1.cards, &dealer.cards} {
			*hand = append(*hand, drawCard())
		}
	}
}

func calculatePointsAllPlayers() {
	player1.points = calculatePoints(player1, -1)
	dealer.points = calculatePoints(dealer, -1)
}

func shouldContinue() bool {
	calculatePointsAllPlayers()
	return !dealer.stand && player1.points < 21
}

func waitPlayersDecision() decision {
	var decided decision
	for true {
		fmt.Printf("Press (h)it or (s)tand: ")
		reader := bufio.NewReader(os.Stdin)
		got, _ := reader.ReadString('\n')
		got = strings.Replace(got, "\n", "", -1)

		if got == "h" {
			fmt.Println("You chose Hit!")
			decided = hit
			break
		} else if got == "s" {
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
		value = min(int(card.Rank), 10)
		if card.Rank == deck.Ace {
			numberAces++
		}
		total = total + value
	}

	// special rule for 21
	if numberAces > 0 && total+10 == 21 {
		total = 21
	}
	return total
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func playersTurn() {
	d := waitPlayersDecision()
	if d == hit {
		player1.cards = append(player1.cards, drawCard())
		player1.points = calculatePoints(player1, -1)
		if player1.points >= 21 {
			player1.stand = true
		}
	} else {
		player1.stand = true
	}
}

func dealersTurn() {
	// drawACard
	if dealer.points <= 16 {
		dealer.cards = append(dealer.cards, drawCard())
		dealer.points = calculatePoints(dealer, -1)
	} else {
		dealer.stand = true
	}
}

func printScreen(endGame bool) {

	// clear screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// show dealer cards only at the end, at the beginning only first card is shown
	var dealer_cards string
	var dealer_points int
	var roundText string
	if endGame {
		dealer_cards = Hand(dealer.cards).String()
		dealer_points = calculatePoints(dealer, -1)
		roundText = "END"
	} else {
		dealer_cards = Hand(dealer.cards[0:1]).String() + ", ** HIDDEN **"
		dealer_points = calculatePoints(dealer, 0)
		roundText = fmt.Sprintf("ROUND %d", gameRound)
	}

	fmt.Printf("[ BLACK JACK  - %s ]\n\n", roundText)
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf(" Dealer (%d cards): %s = %d \n", len(dealer.cards), dealer_cards, dealer_points)
	fmt.Printf(" Player (%d cards): %s = %d \n ", len(player1.cards), Hand(player1.cards), player1.points)
	fmt.Printf("------------------------------------------------------------------------------\n")

	if endGame {
		switch {
		case player1.points > 21:
			fmt.Println("PLAYER BURST!")
		case dealer.points > 21:
			fmt.Println("DEALER BURST!")
		case player1.points > dealer.points:
			fmt.Println("PLAYER WINS!")
		case dealer.points > player1.points:
			fmt.Println("DEALER WINS!")
		case dealer.points == player1.points:
			fmt.Println("DRAW!")
		}
		fmt.Println("Press enter to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
