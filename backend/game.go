package main

import (
	"errors"
	"math/rand"
	"time"
)

// CardType represents the type of card
// In Go, we use constants to define enumerated types
type CardType string

const (
	RedTeam      CardType = "red"
	BlueTeam     CardType = "blue"
	Neutral      CardType = "neutral"
	Assassin     CardType = "assassin"
)

// Card represents a single card on the game board
// struct is Go's way of defining a custom data type with fields
type Card struct {
	Word     string   `json:"word"`     // The word on the card
	Type     CardType `json:"type"`     // The team/type this card belongs to
	Revealed bool     `json:"revealed"` // Whether this card has been revealed
}

// Game represents the entire game state
type Game struct {
	Cards       []Card   `json:"cards"`       // All 25 cards on the board
	CurrentTeam CardType `json:"currentTeam"` // Which team's turn it is
	RedScore    int      `json:"redScore"`    // Number of red cards remaining
	BlueScore   int      `json:"blueScore"`   // Number of blue cards remaining
	GameOver    bool     `json:"gameOver"`    // Is the game finished?
	Winner      string   `json:"winner"`      // Who won? (if game is over)
	ClueWord    string   `json:"clueWord"`       // All 25 cards on the board
	ClueNumber 	int 	 `json:"clueNumber"` // Which team's turn it is

}

// wordList contains all possible words for the game
// In a real game, you'd have a much larger list
var wordList = []string{
	"APPLE", "BANK", "BERLIN", "BOND", "BOOT", "BOTTLE", "BRIDGE",
	"BRUSH", "CAR", "CAT", "CHARGE", "CHINA", "CLOCK", "CLOUD",
	"CLUB", "CODE", "COOK", "COPPER", "CROWN", "CYCLE", "DANCE",
	"DOG", "DRAGON", "DUCK", "EAGLE", "EGYPT", "ENGINE", "FALL",
	"FAN", "FIELD", "FIGHTER", "FIRE", "FISH", "FLUTE", "FORK",
	"FRANCE", "GAME", "GLASS", "GOLD", "GRACE", "GRASS", "GREEN",
	"HAM", "HAND", "HAWK", "HEAD", "HEART", "HELICOPTER", "HIMALAYAS",
	"HOLE", "HOLLYWOOD", "HONEY", "HOOD", "HOOK", "HORN", "HORSE",
	"HOTEL", "ICE", "INDIA", "IRON", "IVORY", "JACK", "JAM",
}

// NewGame creates and initializes a new game
// This is a constructor function (Go doesn't have constructors like OOP languages)
func NewGame() *Game {
	// Seed the random number generator with current time
	// This ensures we get different random results each time
	rand.Seed(time.Now().UnixNano())

	// Create a new game with initial state
	game := &Game{
		Cards:       make([]Card, 25), // Create a slice of 25 cards
		CurrentTeam: RedTeam,           // Red team starts
		RedScore:    9,                 // Red team has 9 cards
		BlueScore:   8,                 // Blue team has 8 cards
		GameOver:    false,
		Winner:      "",
		ClueWord:   "",                 
		ClueNumber:    0,
	}

	// Shuffle and pick 25 random words
	words := pickRandomWords(25)

	// Create the card type distribution
	// 9 red, 8 blue, 7 neutral, 1 assassin = 25 total
	cardTypes := []CardType{}
	
	// Add 9 red cards
	for i := 0; i < 9; i++ {
		cardTypes = append(cardTypes, RedTeam)
	}
	
	// Add 8 blue cards
	for i := 0; i < 8; i++ {
		cardTypes = append(cardTypes, BlueTeam)
	}
	
	// Add 7 neutral cards
	for i := 0; i < 7; i++ {
		cardTypes = append(cardTypes, Neutral)
	}
	
	// Add 1 assassin card
	cardTypes = append(cardTypes, Assassin)

	// Shuffle the card types
	rand.Shuffle(len(cardTypes), func(i, j int) {
		cardTypes[i], cardTypes[j] = cardTypes[j], cardTypes[i]
	})

	// Create cards by combining words with types
	for i := 0; i < 25; i++ {
		game.Cards[i] = Card{
			Word:     words[i],
			Type:     cardTypes[i],
			Revealed: false,
		}
	}

	return game
}

// pickRandomWords selects n random words from the word list
func pickRandomWords(n int) []string {
	// Create a copy of the word list to avoid modifying the original
	words := make([]string, len(wordList))
	copy(words, wordList)

	// Shuffle the words
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	// Return the first n words
	return words[:n]
}

// RevealCard reveals a card at the given index
// This is a method on the Game struct (notice the receiver (g *Game))
func (g *Game) RevealCard(index int) error {
	// Validate the index
	if index < 0 || index >= len(g.Cards) {
		return errors.New("invalid card index")
	}

	// Check if game is already over
	if g.GameOver {
		return errors.New("game is already over")
	}

	card := &g.Cards[index]

	// Check if card is already revealed
	if card.Revealed {
		return errors.New("card already revealed")
	}

	// Reveal the card
	card.Revealed = true

	// Update game state based on card type
	switch card.Type {
	case RedTeam:
		g.RedScore--
		if g.RedScore == 0 {
			g.GameOver = true
			g.Winner = "Red Team"
		}
		// If it's not the red team's card, switch turns
		if g.CurrentTeam != RedTeam {
			g.CurrentTeam = RedTeam
			g.ClueWord = ""
			g.ClueNumber = 0
		}

	case BlueTeam:
		g.BlueScore--
		if g.BlueScore == 0 {
			g.GameOver = true
			g.Winner = "Blue Team"
		}
		// If it's not the blue team's card, switch turns
		if g.CurrentTeam != BlueTeam {
			g.CurrentTeam = BlueTeam
			g.ClueWord = ""
			g.ClueNumber = 0
		}

	case Assassin:
		// Revealing the assassin ends the game
		g.GameOver = true
		// The team that revealed the assassin loses
		if g.CurrentTeam == RedTeam {
			g.Winner = "Blue Team"
		} else {
			g.Winner = "Red Team"
		}

	case Neutral:
		// Neutral card just switches the turn
		if g.CurrentTeam == RedTeam {
			g.CurrentTeam = BlueTeam
			g.ClueWord = ""
			g.ClueNumber = 0
		} else {
			g.CurrentTeam = RedTeam
			g.ClueWord = ""
			g.ClueNumber = 0
		}
	}

	return nil
}

func (g *Game) SubmitClue(word string, number int) error {
	if g.GameOver{
		return errors.New("Game Over!")
	}

	if word == ""{
		return errors.New("No word submitted")
	}

	if number < 1{
		return errors.New("Clue number must be at least 1")
	}

	g.ClueWord = word
	g.ClueNumber = number

	return nil
}
