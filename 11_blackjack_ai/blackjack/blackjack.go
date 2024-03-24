package blackjack

import (
	"bufio"
	"fmt"
	"github.com/felipempda/gophercises/09_deck/deck"
	"os"
	"strings"
)

const (
	_ Move = iota
	MoveStand
	MoveHit
)

type Move int

type Options struct {
	Hands int
	Decks int
}

type GameState struct {
	gameRound int
	gameDeck  deck.Deck
	player1   participant
	dealer    participant
}

func New(ops Options) GameState {
	gs := GameState{}
	gs.gameDeck = deck.New(deck.WithMultipleDecks(atLeast(ops.Decks, 1)), deck.Shuffle)
	// draw 2 cards
	for i := 0; i < 2; i++ {

		for _, p := range []*participant{&gs.player1, &gs.dealer} {
			(*p).cards = append((*p).cards, gs.DrawCard())
			(*p).points = Score((*p).cards...)
		}
	}
	return gs
}

// interface
type AI interface {
	Bet() int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card) string
}

// human implementation of interface
type HumanAI struct {
}

func (ai *HumanAI) Bet() int {
	return 1
}

func (ai *HumanAI) Play(hand []deck.Card, dealer deck.Card) Move {
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

func (ai *HumanAI) Results(player [][]deck.Card, dealer []deck.Card) string {
	dealer_points := Score(dealer...)
	player_points := Score(player[0]...)

	fmt.Printf("[ BLACK JACK  - RESULTS ]\n\n")
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf(" Dealer (%d cards): %s = %d \n", len(dealer), handDeck(dealer), dealer_points)
	fmt.Printf(" Player (%d cards): %s = %d \n ", len(player[0]), handDeck(player[0]), player_points)
	fmt.Printf("------------------------------------------------------------------------------\n")

	// if endGame {
	switch {
	case player_points > 21:
		return "PLAYER BURST!"
	case dealer_points > 21:
		return ("DEALER BURST!")
	case player_points > dealer_points:
		return ("PLAYER WINS!")
	case dealer_points > player_points:
		return ("DEALER WINS!")
	case dealer_points == player_points:
		return ("DRAW!")
	default:
		return ("UNKNOWN!")
	}
}

func Score(cards ...deck.Card) int {
	var total int
	var numberAces int
	for _, card := range cards {
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

func atLeast(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func Soft(cards ...deck.Card) bool {
	var total int
	var numberAces int
	for _, card := range cards {
		var value int
		value = min(int(card.Rank), 10)
		if card.Rank == deck.Ace {
			numberAces++
		}
		total = total + value
	}

	// special rule for 21
	if numberAces > 0 {
		if total+10 == 21 {
			return true
		} else {
			return false
		}
	}
	return false
}

type participant struct {
	cards  handDeck
	points int
	stand  bool
	bet    int
}

func (gs *GameState) DrawCard() deck.Card {
	var card deck.Card
	card, gs.gameDeck = gs.gameDeck[0], gs.gameDeck[1:]
	return card
}

func (gs *GameState) PlayGame(ai AI) string {
	// players Turn
	var move Move
	for gs.gameRound = 1; !gs.player1.stand; gs.gameRound++ {
		move = ai.Play(gs.player1.cards, gs.dealer.cards[0])
		if move == MoveHit {
			gs.player1.cards = append(gs.player1.cards, gs.DrawCard())
			gs.player1.points = Score(gs.player1.cards...)
			if gs.player1.points >= 21 {
				gs.player1.stand = true
			}
		} else if move == MoveStand {
			gs.player1.stand = true
		}
	}

	// Dealers turn
	for ; !gs.dealer.stand; gs.gameRound++ {
		gs.dealersTurn()
	}

	results := ai.Results([][]deck.Card{gs.player1.cards}, gs.dealer.cards)

	return results
}

func (gs *GameState) dealersTurn() {
	if gs.dealer.points <= 16 {
		gs.dealer.cards = append(gs.dealer.cards, gs.DrawCard())
		gs.dealer.points = Score(gs.dealer.cards...)
	} else {
		gs.dealer.stand = true
	}
}

// helper to describe cards
type handDeck []deck.Card

func (h handDeck) String() string {
	strs := make([]string, len(h))
	for i := range strs {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}
