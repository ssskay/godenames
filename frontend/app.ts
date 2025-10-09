// TypeScript uses interfaces to define the shape of objects
// This helps catch errors at compile-time and provides autocomplete

interface Card {
    word: string;
    type: 'red' | 'blue' | 'neutral' | 'assassin';  // Union type - can only be one of these
    revealed: boolean;
}

interface GameState {
    cards: Card[];
    currentTeam: 'red' | 'blue';
    redScore: number;
    blueScore: number;
    gameOver: boolean;
    winner: string;
    clueWord: string;
    clueNumber: number;
}

// Class to manage our game's frontend logic
class CodenamesGame {
    // Private properties (only accessible within this class)
    private apiUrl: string = window.location.hostname === 'localhost' 
    ? 'http://localhost:8080/api'
    : '/api';  // Use relative URL in production
    private spymasterMode: boolean = false;
    private gameState: GameState | null = null;

    // Constructor runs when we create a new instance of the class
    constructor() {
        this.initializeEventListeners();
    }

    // Set up all button click handlers
    private initializeEventListeners(): void {
        // Get references to DOM elements
        const startButton = document.getElementById('start-game-btn') as HTMLButtonElement;
        const newGameButton = document.getElementById('new-game-btn') as HTMLButtonElement;
        const spymasterButton = document.getElementById('spymaster-mode-btn') as HTMLButtonElement;
        const clueButton = document.getElementById('submit-clue-btn') as HTMLButtonElement;

        // Add click event listeners
        // The '?' is optional chaining - it prevents errors if element is null
        startButton?.addEventListener('click', () => this.startNewGame());
        newGameButton?.addEventListener('click', () => this.startNewGame());
        spymasterButton?.addEventListener('click', () => this.toggleSpymasterMode());
        clueButton?.addEventListener('click', () => {
        // Get the input elements
        const wordInput = document.getElementById('clue-word-input') as HTMLInputElement;
        const numberInput = document.getElementById('clue-number-input') as HTMLInputElement;
        
        // Get their values
        const word = wordInput.value.trim();
        const count = parseInt(numberInput.value);
        
        // Validate and submit
        if (word && count >= 1) {
            this.submitClue(word, count);
        } else {
            alert('Please enter a valid clue word and number!');
        }
    });
    }

    private async submitClue(word: string, count: number): Promise<void> {
        try {
            const res = await fetch(`${this.apiUrl}/game/clue`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ word, number: count })
            });
            
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            const gameState: GameState = await res.json();
            this.gameState = gameState;
            this.renderGame(gameState);
            const wordInput = document.getElementById('clue-word-input') as HTMLInputElement;
            const numberInput = document.getElementById('clue-number-input') as HTMLInputElement;
            console.log('word:', word, 'count:', count, 'raw value:', numberInput.value);
            if (wordInput) wordInput.value = '';
            if (numberInput) numberInput.value = '1';
        } catch (err) {
            console.error('Error submitting clue:', err);
            alert('Failed to submit clue. Please try again.');
        }
    }

    // Start a new game by calling the backend API
    private async startNewGame(): Promise<void> {
        try {
            // fetch() is the modern way to make HTTP requests in JavaScript/TypeScript
            // await pauses execution until the promise resolves
            const response = await fetch(`${this.apiUrl}/game/new`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            // Check if the request was successful
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // Parse the JSON response
            const gameState: GameState = await response.json();
            
            // Update our local state and UI
            this.gameState = gameState;
            this.renderGame(gameState);
            
            // Hide game over modal if it's showing
            const gameOverDiv = document.getElementById('game-over');
            gameOverDiv?.classList.add('hidden');

        } catch (error) {
            // TypeScript requires us to handle errors properly
            console.error('Error starting new game:', error);
            alert('Failed to start new game. Make sure the backend server is running!');
        }
    }

    // Render the game board and update all UI elements
    private renderGame(gameState: GameState): void {
        // Update scores
        const redScoreEl = document.getElementById('red-score');
        const blueScoreEl = document.getElementById('blue-score');
        const currentTeamEl = document.getElementById('current-team');

        if (redScoreEl) redScoreEl.textContent = gameState.redScore.toString();
        if (blueScoreEl) blueScoreEl.textContent = gameState.blueScore.toString();
        if (currentTeamEl) {
            currentTeamEl.textContent = gameState.currentTeam === 'red' ? 'Red' : 'Blue';
            // Update text color based on current team
            currentTeamEl.style.color = gameState.currentTeam === 'red' ? '#c62828' : '#1565c0';
        }

            const clueDisplay = document.getElementById('current-clue-display');
            const clueText = document.getElementById('clue-text');
            
            if (clueDisplay && clueText) {
                if (gameState.clueWord && gameState.clueWord !== '') {
                    // There's a clue - show it!
                    clueText.textContent = `${gameState.clueWord.toUpperCase()} - ${gameState.clueNumber}`;
                    clueDisplay.classList.remove('hidden');
                } else {
                    // No clue - hide it
                    clueDisplay.classList.add('hidden');
                }
            }

        // Render the game board
        this.renderBoard(gameState.cards);

        // Check if game is over
        if (gameState.gameOver) {
            this.showGameOver(gameState.winner);
        }
    }

    // Render all cards on the board
    private renderBoard(cards: Card[]): void {
        const boardDiv = document.getElementById('game-board');
        if (!boardDiv) return;

        // Clear existing cards
        boardDiv.innerHTML = '';

        // Create a card element for each card in the game state
        cards.forEach((card, index) => {
            const cardDiv = document.createElement('div');
            cardDiv.className = 'card';
            cardDiv.textContent = card.word;

            // Add revealed class and type if card is revealed
            if (card.revealed) {
                cardDiv.classList.add('revealed', card.type);
            }

            // Add spymaster view class if in spymaster mode (and not revealed)
            if (this.spymasterMode && !card.revealed) {
                cardDiv.classList.add('spymaster-view', card.type);
            }

            // Add click handler to reveal card
            // We use an arrow function to preserve 'this' context
            cardDiv.addEventListener('click', () => this.revealCard(index));

            boardDiv.appendChild(cardDiv);
        });
    }

    // Reveal a card by sending a request to the backend
    private async revealCard(index: number): Promise<void> {
        // Don't allow revealing if game is over or card is already revealed
        if (this.gameState?.gameOver || this.gameState?.cards[index].revealed) {
            return;
        }

        try {
            const response = await fetch(`${this.apiUrl}/game/reveal`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                // Convert JavaScript object to JSON string
                body: JSON.stringify({ index })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // Get updated game state from backend
            const gameState: GameState = await response.json();
            this.gameState = gameState;
            this.renderGame(gameState);

        } catch (error) {
            console.error('Error revealing card:', error);
            alert('Failed to reveal card. Please try again.');
        }
    }

    // Toggle between normal and spymaster view
    private toggleSpymasterMode(): void {
        this.spymasterMode = !this.spymasterMode;
        
        // Re-render the board with the new mode
        if (this.gameState) {
            this.renderBoard(this.gameState.cards);
        }

        // Update button text
        const button = document.getElementById('spymaster-mode-btn');
        if (button) {
            button.textContent = this.spymasterMode 
                ? 'Hide Spymaster View' 
                : 'Show Spymaster View';
        }
    }

    // Show the game over modal
    private showGameOver(winner: string): void {
        const gameOverDiv = document.getElementById('game-over');
        const winnerText = document.getElementById('winner-text');

        if (gameOverDiv && winnerText) {
            winnerText.textContent = `${winner} wins! 🎉`;
            gameOverDiv.classList.remove('hidden');
        }
    }
}

// Create an instance of the game when the page loads
// This is the entry point of our application
const game = new CodenamesGame();

// You could also start a game automatically:
// game.startNewGame();