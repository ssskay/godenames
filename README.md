# Codenames Game - Learning Project

**Learning Goals:** Go backend + TypeScript frontend + REST API communication

---

## 📁 Project Structure

```
codenames/
├── backend/
│   ├── main.go          # HTTP server, API endpoints, routing
│   ├── game.go          # Game logic, structs, methods
│   └── go.mod           # Go dependencies
└── frontend/
    ├── index.html       # UI structure
    ├── styles.css       # Styling
    ├── app.ts           # TypeScript source
    ├── app.js           # Compiled JavaScript (created by tsc)
    ├── tsconfig.json    # TypeScript configuration
    ├── package.json     # Node dependencies
    └── node_modules/    # Dependencies folder
```

---

## ✅ What's Working

### Backend (Go)
- ✅ HTTP server running on port 8080
- ✅ CORS enabled for frontend communication
- ✅ Game state management
- ✅ Word list and board generation (5×5 grid, 25 words)
- ✅ Card type distribution (9 red, 8 blue, 7 neutral, 1 assassin)
- ✅ Card reveal logic
- ✅ Turn switching
- ✅ Win detection
- ✅ Three API endpoints working:
  - `POST /api/game/new` - Create new game
  - `GET /api/game/state` - Get current game state
  - `POST /api/game/reveal` - Reveal a card

### Frontend (TypeScript)
- ✅ TypeScript compilation setup
- ✅ Interactive 5×5 game board
- ✅ Score tracking (Red Team: 9, Blue Team: 8)
- ✅ Current turn indicator
- ✅ Spymaster mode (toggle to see card colors)
- ✅ Card reveal with visual feedback
- ✅ Game over modal with winner announcement
- ✅ Responsive design

### Key Go Concepts Learned
- ✅ Structs and struct tags (`json:"fieldName"`)
- ✅ Methods with receivers `func (g *Game) MethodName()`
- ✅ Pointers (`*Game` vs `Game`)
- ✅ Error handling (returning `error`, checking `if err != nil`)
- ✅ HTTP handlers and routing
- ✅ JSON encoding/decoding
- ✅ Slices and slice operations

### Key TypeScript Concepts Learned
- ✅ Interfaces for type definitions
- ✅ Classes and methods
- ✅ Async/await for API calls
- ✅ Type annotations
- ✅ DOM manipulation
- ✅ Fetch API for HTTP requests
- ✅ Compilation with `tsc`

---

## 🚧 In Progress - Clue/Hint Feature

### What We're Adding
A feature where spymasters can give hints to their team:
- **Clue word**: A single word hint (e.g., "ocean")
- **Clue number**: How many cards relate to the hint (e.g., 3)
- Clue should clear when turn switches

### Backend Progress
✅ Added fields to `Game` struct:
```go
ClueWord   string `json:"clueWord"`
ClueNumber int    `json:"clueNumber"`
```

✅ Created `SubmitClue()` method in `game.go`:
- Validates word is not empty
- Validates number is >= 1
- Validates game is not over
- Updates game state
- Returns error on validation failure

🚧 **Need to Complete:**
1. Create `ClueRequest` struct in `main.go`:
   ```go
   type ClueRequest struct {
       Word   string `json:"word"`
       Number int    `json:"number"`
   }
   ```

2. Create `handleSubmitClue` handler function:
   - Copy pattern from `handleRevealCard`
   - Parse `ClueRequest` instead of `RevealRequest`
   - Call `game.SubmitClue(req.Word, req.Number)`

3. Register route in `main()`:
   ```go
   mux.HandleFunc("/api/game/submitclue", handleSubmitClue)
   ```

4. Add logic to clear clue on turn switch in `RevealCard()` method

### Frontend Progress
🔲 **Not Started Yet** - Will need to:
- Add UI elements (text input + number input)
- Create TypeScript interface for clue
- Add method to call the API
- Display current clue on screen
- Style the clue display area

---

## 🚀 How to Run

### Terminal 1 - Backend
```bash
cd backend
go run main.go game.go
# Wait for: "Server starting on http://localhost:8080"
```

### Terminal 2 - Frontend (if making TypeScript changes)
```bash
cd frontend
npx tsc --watch
# Auto-compiles app.ts → app.js on save
```

### Terminal 3 - Frontend Server
```bash
cd frontend
python3 -m http.server 3000
# Open browser: http://localhost:3000
```

---

## 🧪 Testing Backend with curl

```bash
# Create new game
curl -X POST http://localhost:8080/api/game/new | jq

# Get game state (check for clueWord and clueNumber fields)
curl http://localhost:8080/api/game/state | jq

# Reveal a card
curl -X POST http://localhost:8080/api/game/reveal \
  -H "Content-Type: application/json" \
  -d '{"index":0}'

# Submit a clue (once endpoint is complete)
curl -X POST http://localhost:8080/api/game/submitclue \
  -H "Content-Type: application/json" \
  -d '{"word":"ocean","number":3}'
```

---

## 🐛 Common Issues & Solutions

### TypeScript Won't Compile
**Problem:** `error TS6053: File not found` (parent directory issues)
**Solution:** Use `npx tsc --project ./tsconfig.json`

### Backend Won't Start
**Problem:** Port already in use
**Solution:** Change port in `main.go` or kill process using port 8080

### Frontend Can't Reach Backend
**Problem:** CORS errors or "Failed to fetch"
**Solution:** 
- Make sure backend is running
- Check backend shows "Server starting on http://localhost:8080"
- Verify frontend is using correct URL

### Changes Not Showing
**Problem:** Made changes but nothing updates
**Solution:**
- Backend: Restart server (Ctrl+C, then `go run main.go game.go`)
- Frontend: Recompile TS (`npx tsc`) and refresh browser (Cmd+R)

---

## 📝 Next Steps (After Clue Feature)

### Easy Improvements
1. Add "End Turn" button (manual turn switching)
2. Add guess counter (track how many guesses used)
3. Add animation when cards are revealed
4. Add sound effects
5. Display clue history (list of previous clues)

### Medium Improvements
1. Add timer system for turns
2. Support multiple concurrent games (game IDs)
3. Save/load games (persistence)
4. Add custom word lists
5. Add game statistics tracking

### Advanced Improvements
1. Real-time multiplayer with WebSockets
2. User authentication
3. AI spymaster that generates clues
4. Mobile app version
5. Database integration

---

## 🎓 Key Learning Points

### Go Patterns You Now Know
- **Struct methods:** `func (g *Game) DoSomething() error`
- **Error handling:** `if err != nil { return err }`
- **JSON tags:** `Field string `json:"fieldName"``
- **HTTP handlers:** `func handleX(w http.ResponseWriter, r *http.Request)`
- **Pointers:** When to use `*Type` vs `Type`

### TypeScript Patterns You Now Know
- **Interfaces:** Define object shapes
- **Async/await:** Clean asynchronous code
- **Type safety:** Catch errors at compile time
- **Classes:** Organize related functionality
- **DOM manipulation:** Update UI programmatically

### REST API Design
- **POST:** Create/modify resources (new game, submit clue)
- **GET:** Retrieve resources (game state)
- **Request/Response cycle:** JSON in, JSON out
- **Error handling:** Return appropriate HTTP status codes

---

## 📚 Resources

- **Go Tutorial:** https://go.dev/tour/
- **TypeScript Handbook:** https://www.typescriptlang.org/docs/
- **HTTP Status Codes:** https://httpstatuses.com/
- **JSON & APIs:** https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API

---

## 💪 You've Learned A TON!

From knowing nothing about Go and TypeScript to:
- Building a full HTTP server
- Creating REST APIs
- Managing game state
- Working with JSON
- TypeScript compilation
- Frontend/backend communication
- Error handling in both languages
- Testing with curl
- Debugging techniques

**Keep going! You're doing great!** 🎉
---

Maintained by [Sara Kay](https://sarakay.me) · [@ssskay](https://github.com/ssskay) · [more projects](https://sarakay.me/projects.html)
