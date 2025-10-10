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
	mux := http.NewServeMux()

	// Serve static files from ./frontend (copied by Dockerfile)
	log.Println("Serving static files from: ./frontend")
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	// Register API route handlers
	mux.HandleFunc("/api/game/new", handleNewGame)
	mux.HandleFunc("/api/game/state", handleGameState)
	mux.HandleFunc("/api/game/reveal", handleRevealCard)
	mux.HandleFunc("/api/game/clue", handleSubmitClue)

	// Enable CORS for frontend-backend communication
	handler := enableCORS(mux)

	// Use Railway's PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

// enableCORS wraps our handler to add CORS headers
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	game := NewGame()
	currentGame = game

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

// RevealRequest represents the JSON for revealing a card
type RevealRequest struct {
	Index int `json:"index"`
}

// ClueRequest represents the JSON for submitting a clue
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

	var req RevealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := currentGame.RevealCard(req.Index); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSON(w, currentGame)
}

// handleSubmitClue processes a clue submission
func handleSubmitClue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if currentGame == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	var req ClueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := currentGame.SubmitClue(req.Word, req.Number); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSON(w, currentGame)
}

// sendJSON helper function to send JSON responses
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Global variable to store the current game
var currentGame *Game