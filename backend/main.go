package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// main is the entry point of our Go program
func main() {
	// Create a new HTTP router (multiplexer)
	// This handles routing different URLs to different functions
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("../frontend"))
	mux.Handle("/", fs)

	// Register route handlers
	// When someone visits these URLs, these functions will be called
	mux.HandleFunc("/api/game/new", handleNewGame)
	mux.HandleFunc("/api/game/state", handleGameState)
	mux.HandleFunc("/api/game/reveal", handleRevealCard)
	mux.HandleFunc("/api/game/clue", handleSubmitClue)
	

	// Enable CORS (Cross-Origin Resource Sharing)
	// This allows our frontend (running on a different port) to talk to the backend
	handler := enableCORS(mux)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"  // Default for local development
    }
    
    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(":"+port, handler); err != nil {
        log.Fatal(err)
    }
}

// enableCORS wraps our handler to add CORS headers
// CORS headers tell the browser it's OK for our frontend to access this API
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Remove: origin := r.Header.Get("Origin")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// handleNewGame creates a new game and returns the game state
func handleNewGame(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a new game
	game := NewGame()

	// Save game to our global store (in a real app, use a database)
	currentGame = game

	// Send the game state as JSON
	sendJSON(w, game)
}

// handleGameState returns the current game state
func handleGameState(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	sendJSON(w, currentGame)
}

// RevealRequest represents the JSON we expect when revealing a card
type RevealRequest struct {
	Index int `json:"index"` // The index of the card to reveal
}

type ClueRequest struct {
    Word   string `json:"word"`
    Number int    `json:"number"`
}

// handleRevealCard reveals a card and updates the game state
func handleRevealCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if currentGame == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	// Parse the JSON body
	var req RevealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Reveal the card
	if err := currentGame.RevealCard(req.Index); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated game state
	sendJSON(w, currentGame)
}


func handleSubmitClue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if currentGame == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	// Parse the JSON body
	var req ClueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Reveal the card
	if err := currentGame.SubmitClue(req.Word, req.Number); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated game state
	sendJSON(w, currentGame)
}

// sendJSON is a helper function to send JSON responses
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder writes JSON directly to the response writer
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Global variable to store the current game
// In a real application, you'd use a database or in-memory store like Redis
var currentGame *Game
