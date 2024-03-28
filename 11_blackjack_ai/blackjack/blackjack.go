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
	cards     handDeck
	points    int
	balance   int
	playerBet int
}

type Options struct {
	Hands           int
	Decks           int
	BlackJackPayout float64
}

type GameState struct {
	gameRound       int
	gameDeck        deck.Deck
	player1         participant
	dealer          participant
	state           state
	nDecks          int
	nHands          int
	blackJackPayout float64
}

func New(ops Options) GameState {

	if ops.Decks <= 0 {
		ops.Decks = 3
	}
	if ops.Hands <= 0 {
		ops.Hands = 100
	}
	if ops.BlackJackPayout <= 0.0 {
		ops.BlackJackPayout = 1.5
	}

	return GameState{
		state: statePlayerTurn,
		player1: participant{
			balance: 0,
		},
		dealer: participant{
			balance: 0,
		},
		nDecks:          ops.Decks,
		nHands:          ops.Hands,
		blackJackPayout: ops.BlackJackPayout,
	}
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

	// testing BlackJacK manually:
	// gs.player1.cards = []deck.Card{
	// 	{Rank: deck.Five},
	// 	{Rank: deck.Six},
	// 	{Rank: deck.Ten},
	// }
	// gs.dealer.cards = []deck.Card{
	// 	{Rank: deck.Ace},
	// 	{Rank: deck.Ten},
	// }
	// gs.player1.points = Score(gs.player1.cards...)
	// gs.dealer.points = Score(gs.dealer.cards...)
	gs.state = statePlayerTurn
}

func (gs *GameState) EndGame(ai AI) {
	dealer_points := Score(gs.dealer.cards...)
	player_points := Score(gs.player1.cards...)
	playerBJ, dealerBJ := BlackJack(gs.player1.cards...), BlackJack(gs.dealer.cards...)

	ai.Results([][]deck.Card{gs.player1.cards}, gs.dealer.cards)

	winnings := gs.player1.playerBet
	switch {
	case playerBJ && dealerBJ:
		fmt.Println("BOTH BLACK JACK!")
		winnings = 0
	case dealerBJ:
		fmt.Println("DEALER BLACK JACK!")
		winnings = winnings * -1
	case playerBJ:
		fmt.Println("PLAYER BLACK JACK!")
		winnings = int(float64(gs.player1.playerBet) * gs.blackJackPayout)
	case player_points > 21:
		fmt.Println("PLAYER BURST!")
		winnings = winnings * -1
	case dealer_points > 21:
		fmt.Println("DEALER BURST!")
	case player_points > dealer_points:
		fmt.Println("PLAYER WINS!")
	case dealer_points > player_points:
		fmt.Println("DEALER WINS!")
		winnings = winnings * -1
	case dealer_points == player_points:
		fmt.Println("DRAW!")
		winnings = 0
	default:
		panic("UNKNOWN!")
	}
	fmt.Println()
	gs.player1.cards = nil
	gs.player1.cards = nil
	fmt.Printf(">> Previous Balance: %d\n", gs.player1.balance)
	gs.player1.balance += winnings
	fmt.Printf(">>> Winnings: %d\n", winnings)
	fmt.Printf(">>> New Balance: %d\n\n", gs.player1.balance)

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

	gs.gameDeck = nil
	min := gs.nDecks * 52 / 3

	for n := 1; n <= gs.nHands; n++ {
		shuffled := false
		//fmt.Printf("DECK SIZE: %d \n", len(gs.gameDeck))
		if len(gs.gameDeck) < min {
			//fmt.Print("\n\n\n*****Shuffling deck...******\n\n\n")
			gs.gameDeck = deck.New(deck.WithMultipleDecks(atLeast(gs.nDecks, 1)), deck.Shuffle)
			shuffled = true
		}
		gs.Bet(ai, shuffled)
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

func (gs *GameState) Bet(ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	gs.player1.playerBet = bet
}

func BlackJack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}
