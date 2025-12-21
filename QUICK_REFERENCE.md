# Quick Reference Guide

Fast reference for common tasks and information about the Super Snake codebase.

---

## 📚 Documentation Files

| File | Purpose |
|------|---------|
| [README.md](README.md) | User guide, features, installation |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System design, patterns, data flow |
| [API_REFERENCE.md](API_REFERENCE.md) | Complete API documentation |
| [DEVELOPER_GUIDE.md](DEVELOPER_GUIDE.md) | Development setup, contributing |

---

## 🚀 Quick Start Commands

```bash
# Build
go build -o super_snake

# Run
./super_snake

# Run directly
go run main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...

# View documentation
go doc -all

# Clean build artifacts
go clean
```

---

## 📁 Key Files Location

```
Game Logic       → tui/game/game.go
Menu             → tui/menu/start_game.go
Leaderboard      → tui/leaderboard/leaderboard.go
Scores DB        → internal/score.go
Sessions         → internal/session.go
Database Init    → internal/db.go
Key Input        → utils/keys.go
View Routes      → tui/views/mode.go
Obstacles        → tui/game/pillars.go
```

---

## 🎮 Game Controls

| Action | Keys |
|--------|------|
| Move Up | `W`, `A`, `↑` |
| Move Down | `S`, `B`, `↓` |
| Move Left | `A`, `D`, `←` |
| Move Right | `D`, `C`, `→` |
| Pause | `Space` |
| Confirm/Select | `Enter`, `Space` |
| Back/Exit | `Esc`, `Q` |
| Quit Game | `Ctrl+C`, `Q` |

---

## 💻 Key Data Structures

### Position
```go
type Position struct {
    X int
    Y int
}
```

### Game State
```go
type GameModel struct {
    Snake []Position      // Snake body
    Food Food             // Current food
    Direction Direction   // Current direction
    Score int             // Current score
    IsGameOver bool       // Game state
}
```

### Score Record
```go
type Score struct {
    ID int
    User string
    Session string
    Value int
    CreatedAt time.Time
}
```

---

## 🔧 Common Development Tasks

### Add a New Level

**Files to modify**:
1. `tui/game/pillars.go` - Add obstacle definitions
2. `tui/game/game.go` - Add LevelXGameConfig() function
3. `tui/views/mode.go` - Add ModeGameX constant
4. `tui/super_snake.go` - Add case in NextLevelConfigFromMode()

**Example**:
```go
// In game.go
func Level6GameConfig() GameStartConfig {
    return GameStartConfig{
        Rows: 30, Columns: 60, SpeedMs: 50, Level: 5,
        Pillars: level6Pillars,
        ScoreService: internal.GetScoreService(),
        SessionManager: internal.GetSessionManager(),
    }
}
```

### Change Game Speed

**File**: `tui/game/game.go`
```go
// In LevelXGameConfig() functions, modify SpeedMs
SpeedMs: 150,  // Increase = slower, decrease = faster
```

### Add a New Menu Option

**File**: `tui/menu/start_game.go`
```go
func InitalModel() StartGameModel {
    return StartGameModel{
        choices: []string{"Start Game", "Leaderboard", "NewOption", "Exit"},
        cursor: 0,
    }
}

// In Update(), handle the new option index
if m.cursor == 2 {  // NewOption
    return m, tea.Batch(views.SwitchModeCmd(views.ModeNewView))
}
```

### Query Scores

**File**: Any file
```go
import "github.com/the-Jinxist/golang_snake_game/internal"

scoreService := internal.GetScoreService()
highScore, _ := scoreService.GetHighScore(context.Background())
allScores, _ := scoreService.GetScores(context.Background())
```

### Create New View

**Steps**:
1. Create `tui/newview/model.go`
2. Implement `tea.Model` interface (Init, Update, View)
3. Add mode in `tui/views/mode.go`
4. Add routing in `tui/super_snake.go`

---

## 🐛 Debugging Tips

### Enable Log Output

```go
import "log"

log.SetOutput(os.Stderr)  // Don't interfere with game rendering
log.Printf("Debug: %+v", variable)
```

### Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| Game doesn't start | Check database path, ensure `my.db` writable |
| Input not responding | Verify terminal supports ANSI codes |
| Score not saving | Check database file exists and permissions |
| Collision detection off | Verify Position X/Y calculation |
| Memory leak | Check for unclosed database connections |

---

## 📊 Database Schema

```sql
CREATE TABLE scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user TEXT,                    -- Hostname
    session TEXT UNIQUE,          -- Game session ID
    value INTEGER,                -- Score value
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Common Queries

```sql
-- Get all scores
SELECT * FROM scores ORDER BY value DESC;

-- Get scores for user
SELECT * FROM scores WHERE user = 'username';

-- Get high score
SELECT * FROM scores ORDER BY value DESC LIMIT 1;

-- Get recent scores
SELECT * FROM scores ORDER BY created_at DESC LIMIT 10;
```

---

## 🎯 View Mode Flow

```
SuperSnake (Root)
├── ModeMenu → StartGameModel
├── ModeGame → GameModel (Level 1)
├── ModeGame1-5 → GameModel (Levels 1-5)
├── ModeLeaderboard → Leaderboard
├── ModeGameCompleted → GameModel (Game Over)
└── ModeGameOver → GameModel (Game Over)
```

---

## 🔌 Service Interfaces

### ScoreService
```go
// Get high score
GetHighScore(ctx context.Context) (Score, error)

// Get all scores
GetScores(ctx context.Context) ([]Score, error)

// Save current score
SetCurrentScore(ctx context.Context, value int) error

// Get current game score
GetCurrentScore(ctx context.Context) (int, error)
```

### SessionManager
```go
// Create session
CreateNewSession(value any) (string, error)

// Get current session
GetCurrentSession() (string, error)

// Destroy session
DestroyCurrentSession() error
```

---

## 📦 Dependencies

### Core
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/spf13/cobra` - CLI framework

### Database
- `github.com/mattn/go-sqlite3` - SQLite driver

### Utilities
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/Broderick-Westrope/charmutils` - Utilities

---

## 🎯 Key Functions

### Game Logic

| Function | Purpose |
|----------|---------|
| `moveSnake()` | Update snake position |
| `spawnFood()` | Generate food location |
| `checkCollision()` | Detect collisions |
| `isSnake(x, y)` | Check if position is snake |

### Service

| Function | Purpose |
|----------|---------|
| `GetScoreService()` | Get score service singleton |
| `GetSessionManager()` | Get session manager singleton |
| `IntializeConfigs()` | Initialize all services |

### UI

| Function | Purpose |
|----------|---------|
| `NewModel()` | Create root TUI model |
| `InitalGameModel()` | Create game model |
| `NewLeaderboardModel()` | Create leaderboard |
| `SwitchModeCmd()` | Create mode switch command |

---

## 🎮 Game Difficulty Levels

| Level | Speed | Pillars | Notes |
|-------|-------|---------|-------|
| 1 | 100ms | 1 bar | Introduction |
| 2 | 90ms | 2 obstacles | Moderate |
| 3 | 80ms | 3 obstacles | Challenging |
| 4 | 70ms | 4 obstacles | Hard |
| 5 | 60ms | 5 obstacles | Expert |

(Lower SpeedMs = faster game)

---

## 🔄 Message Flow Example

```
User presses 'W'
    ↓
KeyMsg {String: "w"}
    ↓
SuperSnake.Update(KeyMsg)
    ↓
GameModel.Update(KeyMsg)
    ↓
Parse input → Update Direction
    ↓
moveSnake() → Update Position
    ↓
checkCollision() → Check game state
    ↓
GameModel.View() → Render board
    ↓
Terminal displays update
```

---

## 🧪 Testing Pattern

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal -v

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📝 Commit Message Template

```
type(scope): description

More detailed explanation if needed.

Closes #issue-number
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

---

## 🔗 Useful Links

- [Go Documentation](https://golang.org/doc/)
- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)
- [Project Issues](https://github.com/the-Jinxist/golang_snake_game/issues)
- [SQLite Docs](https://www.sqlite.org/)

---

## 📖 Reading Order (for new developers)

1. **README.md** - Understand what the project does
2. **ARCHITECTURE.md** - Understand how it's structured
3. **API_REFERENCE.md** - Learn the available components
4. **DEVELOPER_GUIDE.md** - Set up development environment
5. **Code** - Read the actual implementation

---

## ✅ Checklist for New Contributions

- [ ] Code follows Go conventions
- [ ] Tests added for new functionality
- [ ] Existing tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] Documentation updated
- [ ] Commit message uses conventional format
- [ ] PR description is clear and complete

---

**Last Updated**: December 2025  
**Maintained by**: The-Jinxist  
**License**: See LICENSE file
