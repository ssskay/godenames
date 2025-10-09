"use strict";
class CodenamesGame {
    constructor() {
        this.apiUrl = 'http://localhost:8080/api';
        this.spymasterMode = false;
        this.gameState = null;
        this.initializeEventListeners();
    }
    initializeEventListeners() {
        const startButton = document.getElementById('start-game-btn');
        const newGameButton = document.getElementById('new-game-btn');
        const spymasterButton = document.getElementById('spymaster-mode-btn');
        const clueButton = document.getElementById('submit-clue-btn');
        startButton?.addEventListener('click', () => this.startNewGame());
        newGameButton?.addEventListener('click', () => this.startNewGame());
        spymasterButton?.addEventListener('click', () => this.toggleSpymasterMode());
        clueButton?.addEventListener('click', () => {
            const wordInput = document.getElementById('clue-word-input');
            const numberInput = document.getElementById('clue-number-input');
            const word = wordInput.value.trim();
            const count = parseInt(numberInput.value);
            if (word && count >= 1) {
                this.submitClue(word, count);
            }
            else {
                alert('Please enter a valid clue word and number!');
            }
        });
    }
    async submitClue(word, count) {
        try {
            const res = await fetch(`${this.apiUrl}/game/clue`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ word, number: count })
            });
            if (!res.ok)
                throw new Error(`HTTP ${res.status}`);
            const gameState = await res.json();
            this.gameState = gameState;
            this.renderGame(gameState);
            const wordInput = document.getElementById('clue-word-input');
            const numberInput = document.getElementById('clue-number-input');
            console.log('word:', word, 'count:', count, 'raw value:', numberInput.value);
            if (wordInput)
                wordInput.value = '';
            if (numberInput)
                numberInput.value = '1';
        }
        catch (err) {
            console.error('Error submitting clue:', err);
            alert('Failed to submit clue. Please try again.');
        }
    }
    async startNewGame() {
        try {
            const response = await fetch(`${this.apiUrl}/game/new`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const gameState = await response.json();
            this.gameState = gameState;
            this.renderGame(gameState);
            const gameOverDiv = document.getElementById('game-over');
            gameOverDiv?.classList.add('hidden');
        }
        catch (error) {
            console.error('Error starting new game:', error);
            alert('Failed to start new game. Make sure the backend server is running!');
        }
    }
    renderGame(gameState) {
        const redScoreEl = document.getElementById('red-score');
        const blueScoreEl = document.getElementById('blue-score');
        const currentTeamEl = document.getElementById('current-team');
        if (redScoreEl)
            redScoreEl.textContent = gameState.redScore.toString();
        if (blueScoreEl)
            blueScoreEl.textContent = gameState.blueScore.toString();
        if (currentTeamEl) {
            currentTeamEl.textContent = gameState.currentTeam === 'red' ? 'Red' : 'Blue';
            currentTeamEl.style.color = gameState.currentTeam === 'red' ? '#c62828' : '#1565c0';
        }
        const clueDisplay = document.getElementById('current-clue-display');
        const clueText = document.getElementById('clue-text');
        if (clueDisplay && clueText) {
            if (gameState.clueWord && gameState.clueWord !== '') {
                clueText.textContent = `${gameState.clueWord.toUpperCase()} - ${gameState.clueNumber}`;
                clueDisplay.classList.remove('hidden');
            }
            else {
                clueDisplay.classList.add('hidden');
            }
        }
        this.renderBoard(gameState.cards);
        if (gameState.gameOver) {
            this.showGameOver(gameState.winner);
        }
    }
    renderBoard(cards) {
        const boardDiv = document.getElementById('game-board');
        if (!boardDiv)
            return;
        boardDiv.innerHTML = '';
        cards.forEach((card, index) => {
            const cardDiv = document.createElement('div');
            cardDiv.className = 'card';
            cardDiv.textContent = card.word;
            if (card.revealed) {
                cardDiv.classList.add('revealed', card.type);
            }
            if (this.spymasterMode && !card.revealed) {
                cardDiv.classList.add('spymaster-view', card.type);
            }
            cardDiv.addEventListener('click', () => this.revealCard(index));
            boardDiv.appendChild(cardDiv);
        });
    }
    async revealCard(index) {
        if (this.gameState?.gameOver || this.gameState?.cards[index].revealed) {
            return;
        }
        try {
            const response = await fetch(`${this.apiUrl}/game/reveal`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ index })
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const gameState = await response.json();
            this.gameState = gameState;
            this.renderGame(gameState);
        }
        catch (error) {
            console.error('Error revealing card:', error);
            alert('Failed to reveal card. Please try again.');
        }
    }
    toggleSpymasterMode() {
        this.spymasterMode = !this.spymasterMode;
        if (this.gameState) {
            this.renderBoard(this.gameState.cards);
        }
        const button = document.getElementById('spymaster-mode-btn');
        if (button) {
            button.textContent = this.spymasterMode
                ? 'Hide Spymaster View'
                : 'Show Spymaster View';
        }
    }
    showGameOver(winner) {
        const gameOverDiv = document.getElementById('game-over');
        const winnerText = document.getElementById('winner-text');
        if (gameOverDiv && winnerText) {
            winnerText.textContent = `${winner} wins! 🎉`;
            gameOverDiv.classList.remove('hidden');
        }
    }
}
const game = new CodenamesGame();
//# sourceMappingURL=app.js.map