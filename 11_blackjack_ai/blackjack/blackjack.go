package blackjack

import (
	"fmt"
	"github.com/felipempda/gophercises/09_deck/deck"
	"strings"
)

const (
	stateBet state = iota
	statePlayerTurn
	stateDealerTurn
	stateHandOver
)

type state int8

type participant struct {
	cards   handDeck
	points  int
	balance int
}

type Options struct {
	Hands int
	Decks int
}

type GameState struct {
	gameRound int
	gameDeck  deck.Deck
	player1   participant
	dealer    participant
	state     state
}

func New(ops Options) GameState {
	gs := GameState{}
	gs.gameDeck = deck.New(deck.WithMultipleDecks(atLeast(ops.Decks, 1)), deck.Shuffle)
	// draw 2 cards
	gs.player1.balance = 0
	gs.state = statePlayerTurn
	return gs
}

func (gs *GameState) InitialDraw() {
	gs.player1.cards = nil
	gs.dealer.cards = nil
	for i := 0; i < 2; i++ {

		for _, p := range []*participant{&gs.player1, &gs.dealer} {
			(*p).cards = append((*p).cards, gs.DrawCard())
			(*p).points = Score((*p).cards...)
		}
	}
	gs.state = statePlayerTurn
}

func (gs *GameState) EndGame(ai AI) {
	dealer_points := Score(gs.dealer.cards...)
	player_points := Score(gs.player1.cards...)

	ai.Results([][]deck.Card{gs.player1.cards}, gs.dealer.cards)
	// if endGame {
	switch {
	case player_points > 21:
		fmt.Println("PLAYER BURST!")
		gs.player1.balance--
	case dealer_points > 21:
		fmt.Println("DEALER BURST!")
		gs.player1.balance++
	case player_points > dealer_points:
		fmt.Println("PLAYER WINS!")
		gs.player1.balance++
	case dealer_points > player_points:
		fmt.Println("DEALER WINS!")
		gs.player1.balance--
	case dealer_points == player_points:
		fmt.Println("DRAW!")
	default:
		fmt.Println("UNKNOWN!")
	}
	fmt.Println()
	gs.player1.cards = nil
	gs.player1.cards = nil
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

func (gs *GameState) DrawCard() deck.Card {
	var card deck.Card
	card, gs.gameDeck = gs.gameDeck[0], gs.gameDeck[1:]
	return card
}

func (gs *GameState) PlayGame(ai AI) int {
	for n := 1; n <= 3; n++ {
		gs.InitialDraw()
		for gs.gameRound = 1; gs.state == statePlayerTurn; gs.gameRound++ {

			// prevent Player to change gameState by passing a copy
			handCopy := make([]deck.Card, len(gs.player1.cards))
			copy(handCopy, gs.player1.cards)

			move := ai.Play(handCopy, gs.dealer.cards[0])
			move(gs)
		}

		for gs.state == stateDealerTurn {
			gs.dealersTurn()
		}

		gs.EndGame(ai)
	}

	return gs.player1.balance
}

func (gs *GameState) dealersTurn() {
	if gs.dealer.points <= 16 || (gs.dealer.points == 17 && Soft(gs.dealer.cards...)) {
		gs.dealer.cards = append(gs.dealer.cards, gs.DrawCard())
		gs.dealer.points = Score(gs.dealer.cards...)
	} else {
		gs.state++
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

type Move func(gs *GameState)

func MoveHit(gs *GameState) {
	participantX := gs.CurrentParticipant()

	participantX.cards = append(participantX.cards, gs.DrawCard())
	participantX.points = Score(participantX.cards...)

	if participantX.points > 21 {
		MoveStand(gs)
	}
}

func MoveStand(gs *GameState) {
	gs.state++
}

func (gs *GameState) CurrentParticipant() *participant {
	switch gs.state {
	case statePlayerTurn:
		return &gs.player1
	case stateDealerTurn:
		return &gs.dealer
	default:
		panic("internal state error")
	}
}
