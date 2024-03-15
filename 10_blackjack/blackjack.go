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

type gameState struct {
	gameRound int
	gameDeck  deck.Deck
	player1   participant
	dealer    participant
}

func (p participant) Clone() participant {
	ret := participant{
		cards:  make([]deck.Card, len(p.cards)),
		points: p.points,
		stand:  p.stand,
	}
	copy(ret.cards, p.cards)
	return ret
}

func (gs gameState) Clone() gameState {
	ret := gameState{
		gameRound: gs.gameRound,
		gameDeck:  make([]deck.Card, len(gs.gameDeck)),
		player1:   gs.player1.Clone(),
		dealer:    gs.dealer.Clone(),
	}
	copy(ret.gameDeck, gs.gameDeck)
	return ret
}

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

func (gs *gameState) drawCard() deck.Card {
	var card deck.Card
	card, gs.gameDeck = gs.gameDeck[0], gs.gameDeck[1:]
	return card
}

func play() {
	for i := 1; i <= 4; i++ {
		gs := initialDraw()
		for gs.gameRound = 1; gs.shouldContinue(); gs.gameRound++ {
			gs.printScreen(false)
			if !gs.player1.stand {
				gs = gs.playersTurn()
			} else {
				gs = gs.dealersTurn()
			}
		}
		gs.printScreen(true)
	}
}

func initialDraw() gameState {
	// start Deck
	gs := gameState{}
	gs.gameDeck = deck.New(deck.WithMultipleDecks(2), deck.Shuffle)
	// draw 2 cards
	for i := 0; i < 2; i++ {

		// array of Hand, so we don't repeat code
		for _, hand := range []*Hand{&gs.player1.cards, &gs.dealer.cards} {
			*hand = append(*hand, gs.drawCard())
		}
	}
	return gs
}

func (gs *gameState) calculatePointsAllPlayers() {
	gs.player1.points = calculatePoints(gs.player1, -1)
	gs.dealer.points = calculatePoints(gs.dealer, -1)
}

func (gs *gameState) shouldContinue() bool {
	gs.calculatePointsAllPlayers()
	return !gs.dealer.stand && gs.player1.points < 21
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

func (gs gameState) playersTurn() gameState {
	ret := gs.Clone()
	d := waitPlayersDecision()
	if d == hit {
		ret.player1.cards = append(ret.player1.cards, ret.drawCard())
		ret.player1.points = calculatePoints(ret.player1, -1)
		if ret.player1.points >= 21 {
			ret.player1.stand = true
		}
	} else {
		ret.player1.stand = true
	}
	return ret
}

func (gs gameState) dealersTurn() gameState {
	// drawACard
	ret := gs.Clone()
	if ret.dealer.points <= 16 {
		ret.dealer.cards = append(ret.dealer.cards, ret.drawCard())
		ret.dealer.points = calculatePoints(ret.dealer, -1)
	} else {
		ret.dealer.stand = true
	}
	return ret
}

func (gs *gameState) printScreen(endGame bool) {

	// clear screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// show dealer cards only at the end, at the beginning only first card is shown
	var dealer_cards string
	var dealer_points int
	var roundText string
	if endGame {
		dealer_cards = Hand(gs.dealer.cards).String()
		dealer_points = calculatePoints(gs.dealer, -1)
		roundText = "END"
	} else {
		dealer_cards = Hand(gs.dealer.cards[0:1]).String() + ", ** HIDDEN **"
		dealer_points = calculatePoints(gs.dealer, 0)
		roundText = fmt.Sprintf("ROUND %d", gs.gameRound)
	}

	fmt.Printf("[ BLACK JACK  - %s ]\n\n", roundText)
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf(" Dealer (%d cards): %s = %d \n", len(gs.dealer.cards), dealer_cards, dealer_points)
	fmt.Printf(" Player (%d cards): %s = %d \n ", len(gs.player1.cards), Hand(gs.player1.cards), gs.player1.points)
	fmt.Printf("------------------------------------------------------------------------------\n")

	if endGame {
		switch {
		case gs.player1.points > 21:
			fmt.Println("PLAYER BURST!")
		case gs.dealer.points > 21:
			fmt.Println("DEALER BURST!")
		case gs.player1.points > gs.dealer.points:
			fmt.Println("PLAYER WINS!")
		case gs.dealer.points > gs.player1.points:
			fmt.Println("DEALER WINS!")
		case gs.dealer.points == gs.player1.points:
			fmt.Println("DRAW!")
		}
		fmt.Println("Press enter to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
