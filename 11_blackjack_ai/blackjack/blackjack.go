package blackjack

import (
	"errors"
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

var (
	errorBust          = errors.New("Player Busted")
	errorCantDouble    = errors.New("Can only double when two cards are in your hand!")
	errorCantSplitSize = errors.New("Can only split when you have two cards!")
	errorCantSplitSame = errors.New("Can only split when cards are be same Rank!")
)

type state int8

type participant struct {
	cards     []handDeck
	cardsIdx  int
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
	dealer          handDeck
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
			cards:   []handDeck{},
			balance: 0,
		},
		dealer:          handDeck{},
		nDecks:          ops.Decks,
		nHands:          ops.Hands,
		blackJackPayout: ops.BlackJackPayout,
	}
}

func (gs *GameState) InitialDraw() {
	gs.player1.cards = make([]handDeck, 1, 1)
	gs.player1.cardsIdx = 0
	gs.dealer = nil
	for i := 0; i < 2; i++ {

		gs.player1.cards[0] = append(gs.player1.cards[0], gs.DrawCard())
		gs.dealer = append(gs.dealer, gs.DrawCard())
	}

	// testing BlackJacK manually:
	// gs.player1.cards[0] = []deck.Card{
	// 	{Rank: deck.Five},
	// 	{Rank: deck.Six},
	// 	{Rank: deck.Ten},
	// }
	// gs.dealer = []deck.Card{
	// 	{Rank: deck.Ace},
	// 	{Rank: deck.Ten},
	// }
	gs.state = statePlayerTurn
}

func (gs *GameState) EndRound(ai AI) {
	dealer_points := Score(gs.dealer...)
	dealerBJ := BlackJack(gs.dealer...)

	allHands := make([][]deck.Card, len(gs.player1.cards))
	for idx, hand := range gs.player1.cards {
		player_points := Score(hand...)
		playerBJ := BlackJack(hand...)
		fmt.Printf("Hand %d\n", idx+1)

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
		default: //case dealer_points == player_points:
			fmt.Println("DRAW!")
			winnings = 0
		}
		fmt.Printf(">> Previous Balance: %d\n", gs.player1.balance)
		gs.player1.balance += winnings
		fmt.Printf(">>> Winnings: %d\n", winnings)
		fmt.Printf(">>> New Balance: %d\n\n", gs.player1.balance)

		allHands[idx] = hand
	}
	ai.Results(allHands, gs.dealer)
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
			currentDeck := gs.CurrentParticipant()
			handCopy := make([]deck.Card, len(gs.player1.cards[gs.player1.cardsIdx]))
			copy(handCopy, *currentDeck)

			move := ai.Play(gs.player1.cardsIdx+1, handCopy, gs.dealer[0])
			err := move(gs)
			switch err {
			case errorBust:
				MoveStand(gs)
			case nil:
				// non-critical errors:
			case errorCantDouble, errorCantSplitSame, errorCantSplitSize:
				fmt.Println(err)
			default:
				panic(err)
			}
		}

		for gs.state == stateDealerTurn {
			gs.dealersTurn()
		}

		gs.EndRound(ai)
	}

	return gs.player1.balance
}

func (gs *GameState) dealersTurn() {
	dealerPoints := Score(gs.dealer...)
	if dealerPoints <= 16 || (dealerPoints == 17 && Soft(gs.dealer...)) {
		gs.dealer = append(gs.dealer, gs.DrawCard())
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

type Move func(gs *GameState) error

func MoveHit(gs *GameState) error {
	currentHand := gs.CurrentParticipant()

	*currentHand = append(*currentHand, gs.DrawCard())

	if Score(*currentHand...) > 21 {
		return errorBust //MoveStand(gs)
	}
	return nil
}

func MoveStand(gs *GameState) error {
	if gs.state == statePlayerTurn {
		gs.player1.cardsIdx++
		if gs.player1.cardsIdx == len(gs.player1.cards) {
			gs.state++
		}
		return nil
	}
	if gs.state == stateDealerTurn {
		gs.state++
	}
	return errors.New("Invalidate State")
}

func MoveDouble(gs *GameState) error {
	total := len(gs.player1.cards)
	if total != 2 {
		return errorCantDouble
	}
	gs.player1.playerBet *= 2
	MoveHit(gs)
	return MoveStand(gs)
}

func (gs *GameState) CurrentParticipant() *handDeck {
	switch gs.state {
	case statePlayerTurn:
		return &gs.player1.cards[gs.player1.cardsIdx]
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

func MoveSplit(gs *GameState) error {
	total := len(gs.player1.cards[gs.player1.cardsIdx])
	if total != 2 {
		return errorCantSplitSize
	}
	if gs.player1.cards[gs.player1.cardsIdx][0].Rank != gs.player1.cards[gs.player1.cardsIdx][1].Rank {
		return errorCantSplitSame
	}

	newDeck := []deck.Card{gs.player1.cards[gs.player1.cardsIdx][1]}
	gs.player1.cards = append(gs.player1.cards, newDeck)
	gs.player1.cards[gs.player1.cardsIdx] = gs.player1.cards[gs.player1.cardsIdx][:1]
	return nil
}
