# Architecture Documentation

## System Design Overview

This document describes the architectural design of the Super Snake application.

---

## 1. Architectural Pattern: Model-View-Update (MVU)

The application implements the **MVU pattern** using the Bubble Tea framework. This pattern ensures:

- **Unidirectional data flow**: Messages trigger updates which produce new models
- **Immutability**: Models are replaced rather than mutated
- **Predictability**: Same input always produces same output
- **Testability**: Pure functions for business logic

### MVU Flow Diagram

```
┌─────────────┐
│   View      │
│  (Render)   │
└──────┬──────┘
       │ displays
       ▼
┌─────────────┐      ┌────────────────┐
│   Model     │──────│   User Input   │
│  (State)    │      │    (KeyMsg)    │
└─────────────┘      └────────────────┘
       ▲                      │
       │                      │ triggers
       └──────────┬───────────┘
                  │
              ┌───▼────┐
              │ Update │
              └────────┘
```

---

## 2. Layered Architecture

The system is organized in distinct layers with clear separation of concerns:

### Layer 1: Presentation Layer (TUI)
**Location**: `tui/`

Components:
- `SuperSnake`: Root model routing
- `game/`: Game rendering and input handling
- `menu/`: Menu interface
- `leaderboard/`: High scores display
- `views/`: View mode definitions

**Responsibilities**:
- Render UI elements
- Handle user input
- Manage view transitions
- Display game state

### Layer 2: Application Logic Layer
**Location**: `tui/game/`, `cmd/`

Components:
- `GameModel`: Game state and mechanics
- `StartGameModel`: Menu logic
- `Leaderboard`: Leaderboard logic

**Responsibilities**:
- Game rules and mechanics
- State transitions
- User interaction logic

### Layer 3: Service Layer
**Location**: `internal/`

Components:
- `ScoreService`: Score operations
- `SessionManager`: Session tracking

**Responsibilities**:
- Business logic abstraction
- Data transformation
- Cross-cutting concerns

### Layer 4: Data Access Layer
**Location**: `internal/db.go`

Components:
- Database initialization
- SQL statement execution
- Connection management

**Responsibilities**:
- Database connectivity
- Schema management
- Query execution

### Layer 5: Infrastructure
**Location**: Root level

Components:
- `main.go`: Entry point
- `cmd/root.go`: CLI initialization

**Responsibilities**:
- Application bootstrap
- Framework initialization
- Configuration loading

---

## 3. Component Interaction Diagram

```
┌──────────────────────────────────────────┐
│           main.go                        │
│        (Program Entry)                   │
└────────────────┬─────────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────────┐
│        cmd/root.go (Cobra)               │
│     - Initialize configs                 │
│     - Launch Bubble Tea                  │
└────────────────┬─────────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────────┐
│      tui/super_snake.go                  │
│    (Root TUI Model - Router)             │
│  ┌─────────────────────────────────────┐ │
│  │ Manages view modes:                 │ │
│  │ - Menu                              │ │
│  │ - Game                              │ │
│  │ - Leaderboard                       │ │
│  │ - Game Over                         │ │
│  └─────────────────────────────────────┘ │
└────────────────┬─────────────────────────┘
                 │
    ┌────────────┼────────────┐
    │            │            │
    ▼            ▼            ▼
┌────────┐   ┌────────┐   ┌──────────┐
│ Menu   │   │ Game   │   │Leaderb.  │
│ Model  │   │ Model  │   │  Model   │
└────────┘   └────────┘   └──────────┘
    │            │            │
    └────────────┼────────────┘
                 │
                 ▼
        ┌─────────────────────┐
        │  internal.go        │
        │  (Services)         │
        ├─────────────────────┤
        │ ScoreService        │
        │ SessionManager      │
        └─────────────────────┘
                 │
                 ▼
        ┌─────────────────────┐
        │  db.go              │
        │  (Database)         │
        ├─────────────────────┤
        │ SQLite Connection   │
        │ Scores Table        │
        └─────────────────────┘
```

---

## 4. Data Flow Patterns

### User Input → Action → State Update → Render

```
User Press 'W'
     │
     ▼
KeyMsg {String: "w"}
     │
     ▼
GameModel.Update(KeyMsg)
     │
     ├─ Parse input
     ├─ Update Direction = Up
     ├─ Move snake
     ├─ Check collisions
     ├─ Update score if food eaten
     │
     ▼
New GameModel (with updated state)
     │
     ▼
GameModel.View() returns rendered string
     │
     ▼
Terminal displays updated board
```

### Score Persistence Flow

```
Game Ends
     │
     ▼
Send SwitchModeMsg(ModeGameCompleted)
     │
     ▼
internal.GetScoreService().SetCurrentScore()
     │
     ▼
INSERT INTO scores (user, session, value)
     │
     ▼
SQLite Database
     │
     ▼
Score Persisted
```

---

## 5. Service Interface Design

### ScoreService Interface

```go
type ScoreService interface {
    // Retrieve the highest score ever recorded
    GetHighScore(ctx context.Context) (Score, error)
    
    // Get all scores (typically ordered by value)
    GetScores(ctx context.Context) ([]Score, error)
    
    // Save the current game's score
    SetCurrentScore(ctx context.Context, value int) error
    
    // Get the current ongoing game score
    GetCurrentScore(ctx context.Context) (int, error)
}
```

**Implementation**: `ScoreServiceImpol`
- Handles database operations
- Tracks current user and session
- Manages score lifecycle

### SessionManager Interface

```go
type SessionManager interface {
    // Create a new unique session identifier
    CreateNewSession(value any) (string, error)
    
    // Destroy the current session
    DestroyCurrentSession() error
    
    // Get the active session ID
    GetCurrentSession() (string, error)
}
```

**Implementation**: `InMemeorySessiomManagerImpl`
- In-memory session storage
- Cryptographically random session IDs
- Session lifecycle management

---

## 6. Dependency Injection

### Initialization Chain

```go
// In cmd/root.go
func init() {
    internal.IntializeConfigs()
}

// In internal/internal.go
func IntializeConfigs() {
    db := CreateDB()                    // Create database
    user, _ := os.Hostname()            // Get system user
    sessionManager = NewSessionManager() // Create session manager
    scoreService = NewScoreService(     // Create score service
        user,
        sessionManager,
        db,
    )
}
```

### Service Injection in Game

```go
// Configuration passed to GameModel
type GameStartConfig struct {
    ScoreService   ScoreService
    SessionManager SessionManager
    Rows           int
    Columns        int
    // ... other config fields
}

// Game uses injected dependencies
gameMod := &GameModel{
    Config: gameConfig,  // Contains all dependencies
    // ... initialize state
}
```

---

## 7. State Management

### Model State Hierarchy

```
SuperSnake (Root Model)
├── width, height (window size)
└── child: tea.Model (current view)
    │
    ├── StartGameModel
    │   ├── choices []string
    │   └── cursor int
    │
    ├── GameModel
    │   ├── Config GameStartConfig
    │   ├── Snake []Position
    │   ├── Food Food
    │   ├── Direction Direction
    │   ├── Score int
    │   ├── IsGameOver bool
    │   ├── IsOutOfBounds bool
    │   ├── spinner spinner.Model
    │   └── isPaused bool
    │
    └── Leaderboard
        ├── Scores []Score
        └── Config LeaderboardConfig
```

### State Transitions

```
START
  │
  ▼
Menu (default)
  │
  ├─[Start Game]──▶ Game Level 1
  │                   │
  │                   ├─[Complete]──▶ GameCompleted
  │                   │                 │
  │                   │                 └─[Continue]──▶ Game Level 2
  │                   │
  │                   └─[Game Over]──▶ Menu
  │
  ├─[Leaderboard]──▶ Leaderboard
  │                   │
  │                   └─[Back]──▶ Menu
  │
  └─[Exit]──▶ EXIT
```

---

## 8. Message Passing Architecture

### Message Types

**1. Keyboard Messages** (Native Bubble Tea)
```go
type KeyMsg struct {
    String() string  // "w", "up", "space", etc.
}
```

**2. Custom Application Messages**
```go
type SwitchModeMsg struct {
    Target Mode  // Which view to switch to
}

type ExitGameMsg struct{}
```

**3. Framework Messages**
```go
tea.ClearScreen()     // Clear terminal
tea.WindowSizeMsg     // Window resized
```

### Message Flow

```
Input Event
     │
     ▼
tea.KeyMsg
     │
     ▼
Model.Update(msg)
     │
     ├─ Handle KeyMsg
     │   └─ Update internal state
     │   └─ Possibly generate new messages
     │
     ├─ Generate SwitchModeMsg if needed
     │   └─ Parent SuperSnake.Update() routes to new child
     │
     └─ Return (updated_model, tea.Cmd)
            │
            └─ tea.Cmd can be:
               ├─ nil (no async action)
               ├─ tea.Quit (exit program)
               ├─ Custom function returning tea.Msg
               └─ tea.Batch(...) (multiple commands)
```

---

## 9. Concurrency & Async Operations

Currently **minimal concurrency**:

- **Synchronous game loop**: Single-threaded main game loop in Update()
- **Database operations**: Blocking SQL queries (could be async with context)
- **Message handling**: Sequential message processing

**Future optimization opportunities**:
1. Async database operations with goroutines
2. Concurrent rendering for complex boards
3. Background leaderboard refresh

---

## 10. Error Handling Strategy

### Error Handling Layers

```
Layer 1: Initialization (fatal errors)
  ├─ Database creation
  ├─ Config loading
  └─ Exit on failure

Layer 2: Service Operations (recoverable errors)
  ├─ Database queries
  ├─ Score persistence
  └─ Log errors, use defaults

Layer 3: User Input (validation)
  ├─ Parse keyboard input
  ├─ Validate game state
  └─ Gracefully handle invalid input

Layer 4: UI Rendering (never fail)
  ├─ Always produce view output
  └─ Fallback rendering available
```

### Error Recovery

```go
// Graceful degradation example
score, err := scoreService.GetCurrentScore(ctx)
if err != nil {
    log.Printf("Error getting score: %v", err)
    score = 0  // Default to 0
}
```

---

## 11. Performance Considerations

### Current Performance

- **Game tick rate**: ~60 FPS (Bubble Tea default)
- **Database queries**: On-demand (not cached)
- **Memory**: Minimal (snake body array, game state)
- **CPU**: Low (simple collision detection)

### Optimization Points

1. **Database Caching**: Cache leaderboard in memory
2. **Collision Detection**: Use hash maps instead of iteration
3. **Rendering**: Only render changed regions
4. **Snake Body**: Use deque instead of slice for efficient removal

---

## 12. Security Considerations

### Current Security

- **Session IDs**: Cryptographically random (20 bytes base64)
- **Database**: Local SQLite (no network exposure)
- **Input Validation**: Terminal input only, limited attack surface

### Future Considerations

- Input sanitization for network multiplayer
- Database encryption at rest
- Session token expiration
- User authentication if adding server component

---

## 13. Testing Strategy

### Testable Components

```go
// Pure functions (easily testable)
func (g *GameModel) isSnake(x, y int) bool { ... }
func (g *GameModel) checkCollision() bool { ... }
func KeyMatchesInput(input string, keys ...Key) bool { ... }

// Service methods (need mocks/integration tests)
func (s *ScoreServiceImpol) GetHighScore(ctx context.Context) (Score, error) { ... }

// UI Models (integration tests with Bubble Tea helpers)
func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
```

### Recommended Test Coverage

1. **Unit Tests**: Game logic, utilities
2. **Integration Tests**: Database operations, services
3. **UI Tests**: Model update flows, state transitions
4. **End-to-End Tests**: Full game scenarios

---

## 14. Deployment Architecture

### Current Deployment

```
Source Code (GitHub)
     │
     ▼
Go Build
     │
     ▼
Binary Executable (super_snake)
     │
     ▼
Local Machine (macOS/Linux/Windows)
     │
     ▼
my.db (SQLite database, local)
```

### Future Deployment Options

1. **Container**: Docker image for consistent environment
2. **Package Manager**: Homebrew, apt, etc.
3. **Web Version**: WASM + Browser UI
4. **Cloud**: Multiplayer server component

---

## 15. Configuration Management

### Configuration Sources

```
Application Configuration
├── Hardcoded Defaults
│   ├── Game board size (Rows, Columns)
│   ├── Game tick speed
│   └── Level configurations
│
├── Runtime Determined
│   ├── System hostname (for user tracking)
│   ├── Terminal size
│   └── Database file location
│
└── User Input
    └── Level selection
    └── Game mode choice
```

### Game Configurations

```go
type GameStartConfig struct {
    ScoreService   ScoreService
    SessionManager SessionManager
    Rows           int
    Columns        int
    SpeedMs        int      // Tick speed in milliseconds
    Level          int      // Difficulty level
    Pillars        []Position  // Obstacle positions
}

// Predefined configs for each level
func Level1GameConfig() GameStartConfig { ... }
func Level2GameConfig() GameStartConfig { ... }
// ... up to Level5GameConfig()
```

---

## 16. Extension Points

The architecture allows for easy extension:

### Add New Game Mode

1. Create new type implementing `tea.Model`
2. Add new `Mode` constant in `views/mode.go`
3. Add case in `SuperSnake.setChild()`
4. Handle routing in parent Update()

### Add New Service

1. Define interface in `internal/`
2. Implement concrete type
3. Initialize in `internal.IntializeConfigs()`
4. Inject via configuration structs
5. Use in models

### Add New View

1. Create model in `tui/`
2. Implement `tea.Model` interface
3. Define mode and routing
4. Handle state transitions

---

## Conclusion

The Super Snake architecture provides:

✅ Clear separation of concerns  
✅ Easy to test and extend  
✅ Scalable component design  
✅ Follows Go best practices  
✅ Leverages proven Bubble Tea patterns  

The MVU + layered approach ensures the codebase remains maintainable as new features are added.
