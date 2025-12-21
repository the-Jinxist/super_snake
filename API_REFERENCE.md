# API Reference Documentation

Complete reference for all public types, interfaces, and functions in the Super Snake codebase.

---

## Table of Contents

- [cmd Package](#cmd-package)
- [internal Package](#internal-package)
- [tui Package](#tui-package)
- [tui/game Package](#tuigame-package)
- [tui/menu Package](#tuimenu-package)
- [tui/leaderboard Package](#tuileaderboard-package)
- [tui/views Package](#tuiviews-package)
- [utils Package](#utils-package)

---

## cmd Package

### Functions

#### Execute()
```go
func Execute()
```

**Description**: Executes the root Cobra command. This is the main entry point for the CLI application.

**Returns**: void (exits program on error)

**Usage**:
```go
package main

func main() {
    cmd.Execute()  // Starts the game
}
```

---

## internal Package

### Types

#### Score
```go
type Score struct {
    ID        int       `db:"id"`
    User      string    `db:"user"`
    Session   string    `db:"session"`
    Value     int       `db:"value"`
    CreatedAt time.Time `db:"created_at"`
}
```

**Description**: Represents a game score record in the database.

**Fields**:
- `ID`: Unique identifier for the score record
- `User`: System hostname/username of the player
- `Session`: Unique session identifier for the game
- `Value`: The final score achieved
- `CreatedAt`: Timestamp of when the score was recorded

---

### Interfaces

#### ScoreService
```go
type ScoreService interface {
    GetHighScore(ctx context.Context) (Score, error)
    GetScores(ctx context.Context) ([]Score, error)
    SetCurrentScore(ctx context.Context, value int) error
    GetCurrentScore(ctx context.Context) (int, error)
}
```

**Description**: Interface for score management operations.

**Methods**:

##### GetHighScore(ctx context.Context) (Score, error)
- **Description**: Retrieves the highest score ever recorded
- **Parameters**:
  - `ctx`: Context for cancellation/timeout
- **Returns**: 
  - `Score`: The highest score record
  - `error`: Non-nil if operation fails
- **Example**:
```go
score, err := scoreService.GetHighScore(context.Background())
if err != nil {
    log.Printf("Error: %v", err)
    return
}
fmt.Printf("High score: %d by %s\n", score.Value, score.User)
```

##### GetScores(ctx context.Context) ([]Score, error)
- **Description**: Retrieves all recorded scores (typically ordered by value)
- **Parameters**:
  - `ctx`: Context for cancellation/timeout
- **Returns**: 
  - `[]Score`: Slice of score records
  - `error`: Non-nil if operation fails
- **Example**:
```go
scores, err := scoreService.GetScores(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, score := range scores {
    fmt.Printf("%s: %d\n", score.User, score.Value)
}
```

##### SetCurrentScore(ctx context.Context, value int) error
- **Description**: Saves the current game's final score
- **Parameters**:
  - `ctx`: Context for cancellation/timeout
  - `value`: The score value to save
- **Returns**: 
  - `error`: Non-nil if save operation fails
- **Example**:
```go
err := scoreService.SetCurrentScore(context.Background(), 150)
if err != nil {
    log.Printf("Failed to save score: %v", err)
}
```

##### GetCurrentScore(ctx context.Context) (int, error)
- **Description**: Gets the score for the current ongoing game session
- **Parameters**:
  - `ctx`: Context for cancellation/timeout
- **Returns**: 
  - `int`: Current game score
  - `error`: Non-nil if operation fails
- **Example**:
```go
score, err := scoreService.GetCurrentScore(context.Background())
if err != nil {
    score = 0  // Default to 0
}
```

---

#### SessionManager
```go
type SessionManager interface {
    CreateNewSession(value any) (string, error)
    DestroyCurrentSession() error
    GetCurrentSession() (string, error)
}
```

**Description**: Interface for managing user game sessions.

**Methods**:

##### CreateNewSession(value any) (string, error)
- **Description**: Creates a new unique session identifier
- **Parameters**:
  - `value`: Optional value to associate with session (currently unused)
- **Returns**:
  - `string`: New session ID (20-character base64)
  - `error`: Non-nil if creation fails
- **Example**:
```go
sessionID, err := sessionManager.CreateNewSession(nil)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("New session: %s\n", sessionID)
```

##### DestroyCurrentSession() error
- **Description**: Clears the current active session
- **Returns**:
  - `error`: Non-nil if operation fails
- **Example**:
```go
err := sessionManager.DestroyCurrentSession()
if err != nil {
    log.Printf("Error destroying session: %v", err)
}
```

##### GetCurrentSession() (string, error)
- **Description**: Retrieves the current active session ID
- **Returns**:
  - `string`: Current session ID
  - `error`: Non-nil if no session exists
- **Example**:
```go
sessionID, err := sessionManager.GetCurrentSession()
if err != nil {
    log.Println("No active session")
    return
}
fmt.Printf("Current session: %s\n", sessionID)
```

---

### Functions

#### GetSessionManager() SessionManager
```go
func GetSessionManager() SessionManager
```

**Description**: Returns the global session manager instance (singleton).

**Returns**: Initialized `SessionManager` implementation

**Usage**:
```go
sessionMgr := internal.GetSessionManager()
sessionID, _ := sessionMgr.GetCurrentSession()
```

---

#### GetScoreService() ScoreService
```go
func GetScoreService() ScoreService
```

**Description**: Returns the global score service instance (singleton).

**Returns**: Initialized `ScoreService` implementation

**Usage**:
```go
scoreService := internal.GetScoreService()
highScore, _ := scoreService.GetHighScore(context.Background())
```

---

#### IntializeConfigs()
```go
func IntializeConfigs()
```

**Description**: Initializes all global configuration and services. Must be called once at startup.

**Responsibilities**:
- Creates SQLite database
- Retrieves system hostname
- Initializes session manager
- Initializes score service

**Panics**: Fatally exits if database cannot be created

**Usage**:
```go
func init() {
    internal.IntializeConfigs()
}
```

---

#### CreateDB() *sql.DB
```go
func CreateDB() *sql.DB
```

**Description**: Creates and initializes the SQLite database connection and schema.

**Returns**: 
- `*sql.DB`: Database connection (or nil on error)

**Side Effects**:
- Creates `my.db` file in current directory
- Creates `scores` table if not exists
- Logs errors to stdout

**Usage**:
```go
db := internal.CreateDB()
if db == nil {
    log.Fatal("Database initialization failed")
}
```

---

#### NewScoreService(user string, sessionMgr SessionManager, db *sql.DB) ScoreService
```go
func NewScoreService(user string, sessionMgr SessionManager, db *sql.DB) ScoreService
```

**Description**: Creates a new score service instance.

**Parameters**:
- `user`: Username/hostname for score tracking
- `sessionMgr`: Session manager for session tracking
- `db`: Database connection for persistence

**Returns**: `ScoreService` implementation

**Usage**:
```go
scoreService := internal.NewScoreService(
    "john",
    sessionManager,
    database,
)
```

---

#### NewSessionManager() SessionManager
```go
func NewSessionManager() SessionManager
```

**Description**: Creates a new in-memory session manager instance.

**Returns**: `SessionManager` implementation

**Usage**:
```go
sessionMgr := internal.NewSessionManager()
```

---

## tui Package

### Types

#### SuperSnake
```go
type SuperSnake struct {
    child  tea.Model  // Current child view
    width  int        // Terminal width
    height int        // Terminal height
}
```

**Description**: Root TUI model that manages view transitions between menu, game, and leaderboard.

**Methods**:

##### NewModel() *SuperSnake
```go
func NewModel() *SuperSnake
```

**Description**: Creates the root TUI model with initial menu view.

**Returns**: `*SuperSnake` instance

**Usage**:
```go
program := tea.NewProgram(tui.NewModel())
program.Run()
```

---

##### Init() tea.Cmd
```go
func (s *SuperSnake) Init() tea.Cmd
```

**Description**: Initializes the model (Bubble Tea interface).

**Returns**: `tea.Cmd` (nil or command to run)

---

##### Update(msg tea.Msg) (tea.Model, tea.Cmd)
```go
func (s *SuperSnake) Update(msg tea.Msg) (tea.Model, tea.Cmd)
```

**Description**: Processes messages and updates model state (Bubble Tea interface).

**Parameters**:
- `msg`: Message to process (`KeyMsg`, `SwitchModeMsg`, etc.)

**Returns**: 
- Updated model
- Command to execute

**Handled Messages**:
- `tea.KeyMsg`: Handles quit commands (Ctrl+C, Q)
- `views.SwitchModeMsg`: Routes to appropriate child view
- Other messages: Delegated to child model

**Example**:
```go
model, cmd := superSnake.Update(tea.KeyMsg{String: "q"})
// Returns model with tea.Quit command
```

---

##### View() string
```go
func (s *SuperSnake) View() string
```

**Description**: Renders the current view (Bubble Tea interface).

**Returns**: String representation of current view

---

### Functions

#### NextLevelConfigFromMode(level views.Mode) game.GameStartConfig
```go
func NextLevelConfigFromMode(level views.Mode) game.GameStartConfig
```

**Description**: Maps view mode to game configuration for the corresponding level.

**Parameters**:
- `level`: View mode (ModeGame1, ModeGame2, etc.)

**Returns**: `GameStartConfig` for that level

**Mapping**:
- `ModeGame1` → Level 1 config
- `ModeGame2` → Level 2 config
- ... up to Level 5
- Default → Level 1 config

---

## tui/game Package

### Types

#### Direction
```go
type Direction int

const (
    Up Direction = iota
    Down
    Left
    Right
)
```

**Description**: Enumeration for snake movement directions.

---

#### Position
```go
type Position struct {
    X int  // Horizontal coordinate
    Y int  // Vertical coordinate
}
```

**Description**: Represents a point on the game board.

---

#### Food
```go
type Food struct {
    Position Position  // Location of food
    BigFish  bool      // Whether it's special food
}
```

**Description**: Represents food on the game board.

**Fields**:
- `Position`: X, Y coordinates
- `BigFish`: If true, grants extra score

---

#### GameModel
```go
type GameModel struct {
    Config       GameStartConfig
    Snake        []Position
    Food         Food
    Direction    Direction
    Score        int
    IsGameOver   bool
    IsOutOfBounds bool
    spinner      spinner.Model
    isPaused     bool
}
```

**Description**: Main game model containing all game state.

**Methods**:

##### InitalGameModel(gameConfig GameStartConfig) *GameModel
```go
func InitalGameModel(gameConfig GameStartConfig) *GameModel
```

**Description**: Creates a new game model with given configuration.

**Parameters**:
- `gameConfig`: Game configuration (speed, board size, etc.)

**Returns**: Initialized `*GameModel`

**Initialization**:
- Snake positioned at center
- Direction set to Right
- Score retrieved from score service
- Spinner initialized for UI

**Example**:
```go
config := game.DefaultGameConfig()
gameModel := game.InitalGameModel(config)
```

---

##### Init() tea.Cmd
```go
func (g *GameModel) Init() tea.Cmd
```

**Description**: Initializes the game (Bubble Tea interface).

**Returns**: Initialization command

---

##### Update(msg tea.Msg) (tea.Model, tea.Cmd)
```go
func (g *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
```

**Description**: Processes game updates (Bubble Tea interface).

**Handles**:
- Keyboard input (direction changes)
- Game timer ticks
- State transitions

**Game Loop**:
1. Process input
2. Move snake
3. Check collisions
4. Regenerate food if eaten
5. Update score
6. Check game over conditions

---

##### View() string
```go
func (g *GameModel) View() string
```

**Description**: Renders the game board (Bubble Tea interface).

**Returns**: Styled ASCII representation of the game

---

#### GameStartConfig
```go
type GameStartConfig struct {
    ScoreService   internal.ScoreService
    SessionManager internal.SessionManager
    Rows           int
    Columns        int
    SpeedMs        int
    Level          int
    Pillars        []Position
}
```

**Description**: Configuration for game initialization.

**Fields**:
- `ScoreService`: For score operations
- `SessionManager`: For session tracking
- `Rows`: Game board height
- `Columns`: Game board width
- `SpeedMs`: Tick speed in milliseconds
- `Level`: Difficulty level
- `Pillars`: Obstacle positions

---

### Functions

#### DefaultGameConfig() GameStartConfig
```go
func DefaultGameConfig() GameStartConfig
```

**Description**: Returns default game configuration (Level 1).

**Returns**: Configured `GameStartConfig`

---

#### Level1GameConfig() GameStartConfig
```go
func Level1GameConfig() GameStartConfig
```

**Description**: Returns configuration for difficulty level 1.

**Returns**: Level 1 `GameStartConfig`

---

#### Level2GameConfig() GameStartConfig, Level3GameConfig() GameStartConfig, Level4GameConfig() GameStartConfig, Level5GameConfig() GameStartConfig

**Description**: Similar to Level1, returns configuration for levels 2-5.

---

## tui/menu Package

### Types

#### StartGameModel
```go
type StartGameModel struct {
    choices []string  // Menu options
    cursor  int       // Selected item index
}
```

**Description**: Model for the main menu screen.

**Methods**:

##### InitalModel() StartGameModel
```go
func InitalModel() StartGameModel
```

**Description**: Creates the initial menu model.

**Returns**: `StartGameModel` with default menu items

---

##### Init() tea.Cmd
```go
func (m StartGameModel) Init() tea.Cmd
```

**Description**: Initializes the menu (Bubble Tea interface).

---

##### Update(msg tea.Msg) (tea.Model, tea.Cmd)
```go
func (m StartGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
```

**Description**: Handles menu input (Bubble Tea interface).

**Handles**:
- Arrow up/W: Move cursor up
- Arrow down/S: Move cursor down
- Enter/Space: Select option

**Options**:
1. Start Game → `views.ModeGame`
2. Leaderboard → `views.ModeLeaderboard`
3. Exit → `tea.Quit`

---

##### View() string
```go
func (m StartGameModel) View() string
```

**Description**: Renders the menu screen (Bubble Tea interface).

**Includes**:
- ASCII art title
- Current high score
- Menu items with cursor
- Help text

---

## tui/leaderboard Package

### Types

#### Leaderboard
```go
type Leaderboard struct {
    Scores []internal.Score
    Config LeaderboardConfig
}
```

**Description**: Model for displaying the leaderboard.

**Methods**:

##### NewLeaderboardModel(config LeaderboardConfig) *Leaderboard
```go
func NewLeaderboardModel(config LeaderboardConfig) *Leaderboard
```

**Description**: Creates a new leaderboard model.

**Parameters**:
- `config`: Leaderboard configuration

**Returns**: `*Leaderboard` with scores loaded from database

---

##### Init() tea.Cmd
```go
func (l *Leaderboard) Init() tea.Cmd
```

**Description**: Initializes the leaderboard (Bubble Tea interface).

---

##### Update(msg tea.Msg) (tea.Model, tea.Cmd)
```go
func (l *Leaderboard) Update(msg tea.Msg) (tea.Model, tea.Cmd)
```

**Description**: Handles leaderboard input (Bubble Tea interface).

**Handles**:
- Escape/Q: Return to menu
- Space: Return to menu

---

##### View() string
```go
func (l *Leaderboard) View() string
```

**Description**: Renders the leaderboard (Bubble Tea interface).

**Displays**:
- ASCII art title "LEADERBOARD"
- Top scores in table format
- User, score value, and date created

---

#### LeaderboardConfig
```go
type LeaderboardConfig struct {
    ScoreService internal.ScoreService
}
```

**Description**: Configuration for leaderboard.

---

### Functions

#### DefaultLeaderboardConfig() LeaderboardConfig
```go
func DefaultLeaderboardConfig() LeaderboardConfig
```

**Description**: Returns default leaderboard configuration.

---

## tui/views Package

### Types

#### Mode
```go
type Mode int

const (
    ModeMenu Mode = iota
    ModeGame
    ModeGame1
    ModeGame2
    ModeGame3
    ModeGame4
    ModeGame5
    ModeLeaderboard
    ModeGameOver
    ModeGameCompleted
)
```

**Description**: Enumeration of all possible view modes.

---

#### SwitchModeMsg
```go
type SwitchModeMsg struct {
    Target Mode  // Mode to switch to
}
```

**Description**: Message to trigger view mode change.

---

#### ExitGameMsg
```go
type ExitGameMsg struct{}
```

**Description**: Message to signal game exit.

---

### Functions

#### NextLevelModeFromCurrent(level int) Mode
```go
func NextLevelModeFromCurrent(level int) Mode
```

**Description**: Maps level number to corresponding mode.

**Parameters**:
- `level`: Level number (0-4)

**Returns**: Corresponding `Mode` (ModeGame1-5)

**Mapping**:
- 0 → ModeGame1
- 1 → ModeGame2
- 2 → ModeGame3
- 3 → ModeGame4
- 4 → ModeGame5

---

#### ClearScreen() tea.Cmd
```go
func ClearScreen() tea.Cmd
```

**Description**: Returns a command to clear the terminal screen.

**Returns**: `tea.Cmd` that clears screen

---

#### SwitchModeCmd(target Mode) tea.Cmd
```go
func SwitchModeCmd(target Mode) tea.Cmd
```

**Description**: Returns a command to switch to specified mode.

**Parameters**:
- `target`: Mode to switch to

**Returns**: `tea.Cmd` that generates `SwitchModeMsg`

---

## utils Package

### Types

#### Key
```go
type Key []string
```

**Description**: Type representing multiple key input variations.

---

### Variables

```go
var (
    KeyUp    Key = []string{"A", "w", "up"}
    KeyDown  Key = []string{"B", "s", "down"}
    KeyRight Key = []string{"C", "d", "right"}
    KeyLeft  Key = []string{"D", "a", "left"}
    Enter    Key = []string{"enter"}
    Esc      Key = []string{"esc"}
    Space    Key = []string{" "}
)
```

**Description**: Predefined key bindings supporting multiple input formats.

**ANSI Terminal Codes**:
- "A" = Up arrow (ANSI)
- "B" = Down arrow (ANSI)
- "C" = Right arrow (ANSI)
- "D" = Left arrow (ANSI)

---

### Functions

#### KeyMatchesInput(input string, keys ...Key) bool
```go
func KeyMatchesInput(input string, keys ...Key) bool
```

**Description**: Checks if input matches any of the provided key definitions.

**Parameters**:
- `input`: User input string
- `keys`: Variable number of `Key` definitions to check against

**Returns**: 
- `true` if input matches any key definition
- `false` otherwise

**Example**:
```go
if utils.KeyMatchesInput("w", utils.KeyUp) {
    fmt.Println("User pressed up!")
}

if utils.KeyMatchesInput("A", utils.KeyUp, utils.KeyDown) {
    fmt.Println("User pressed up or down!")
}
```

---

## Appendix: Common Usage Patterns

### Starting a Game

```go
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/the-Jinxist/golang_snake_game/cmd"
    "github.com/the-Jinxist/golang_snake_game/tui"
)

func main() {
    // Initialize configs (database, services)
    cmd.Execute()  // This calls internal.IntializeConfigs() first
}
```

### Accessing Score Service

```go
import (
    "context"
    "github.com/the-Jinxist/golang_snake_game/internal"
)

func showHighScore() {
    scoreService := internal.GetScoreService()
    score, err := scoreService.GetHighScore(context.Background())
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    fmt.Printf("High score: %d\n", score.Value)
}
```

### Custom Game Configuration

```go
import "github.com/the-Jinxist/golang_snake_game/tui/game"

// Create custom level with specific settings
config := game.GameStartConfig{
    ScoreService: internal.GetScoreService(),
    SessionManager: internal.GetSessionManager(),
    Rows: 25,
    Columns: 50,
    SpeedMs: 100,
    Level: 3,
    Pillars: myCustomPillars,
}

gameModel := game.InitalGameModel(config)
```

---

**Last Updated**: December 2025  
**Version**: 1.0  
**Go Version**: 1.24.2+
