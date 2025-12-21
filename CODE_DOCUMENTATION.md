# Code Documentation Guidelines

Guidelines for maintaining comprehensive documentation within the codebase through comments and docstrings.

---

## Table of Contents

- [Package Documentation](#package-documentation)
- [Function Documentation](#function-documentation)
- [Type Documentation](#type-documentation)
- [Complex Logic Documentation](#complex-logic-documentation)
- [Examples](#examples)
- [Best Practices](#best-practices)

---

## Package Documentation

Every package should have a package-level comment explaining its purpose.

### Format

**Location**: First line of package file (usually in a `doc.go` file or top of main package file)

```go
// Package game contains the core snake game logic.
//
// The GameModel type implements the tea.Model interface and handles:
// - Snake movement and collision detection
// - Food spawning and consumption
// - Score tracking
// - Game state management (playing, paused, game over)
//
// Game configurations can be created for different difficulty levels using
// the provided Level*GameConfig() functions.
package game
```

### What to Include

- Brief one-sentence description
- What the package does
- Key types/functions
- Common usage patterns
- Any package-level setup requirements

---

## Function Documentation

Every public function should have a documentation comment.

### Format

```go
// FunctionName does something specific.
//
// More detailed explanation if the function's purpose isn't immediately clear.
// Can span multiple lines explaining behavior, edge cases, or complex logic.
//
// Parameters:
//   - param1: Description of what param1 does
//   - param2: Description of what param2 does
//
// Returns:
//   - ResultType: What is returned and when
//   - error: Error conditions (if applicable)
//
// Example:
//   result, err := FunctionName(value)
//   if err != nil {
//       // handle error
//   }
func FunctionName(param1 string, param2 int) (string, error) {
    // implementation
}
```

### Specific Patterns

#### For Methods on Types

```go
// Update processes input messages and updates game state.
//
// It handles keyboard input to change snake direction, timer ticks for
// movement, and generates rendering updates. Returns new model and any
// commands to execute.
func (g *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // implementation
}
```

#### For Interface Methods

```go
// GetHighScore retrieves the highest score ever recorded.
//
// The returned Score contains the user, session ID, score value, and
// when it was recorded. Returns an error if the database cannot be accessed.
//
// Returns:
//   - Score: The highest score record (ordered by Value field)
//   - error: Non-nil if database query fails
func (s *ScoreServiceImpl) GetHighScore(ctx context.Context) (Score, error) {
    // implementation
}
```

#### For Constructors

```go
// NewLeaderboardModel creates a new leaderboard display model.
//
// Loads all scores from the database and initializes the leaderboard view.
// Scores are automatically ordered by value in descending order.
//
// Parameters:
//   - config: Leaderboard configuration (must contain valid ScoreService)
//
// Returns:
//   - *Leaderboard: Initialized leaderboard ready to use
func NewLeaderboardModel(config LeaderboardConfig) *Leaderboard {
    // implementation
}
```

---

## Type Documentation

Every public type should have a documentation comment.

### Format

```go
// TypeName describes what this type represents.
//
// More detailed explanation of the purpose, usage, and constraints.
// Explain invariants that must be maintained.
//
// Fields:
//   - FieldA: Description of FieldA and valid values
//   - FieldB: Description of FieldB and valid values
type TypeName struct {
    FieldA string
    FieldB int
}
```

### Specific Patterns

#### For Structs

```go
// GameModel represents the complete state of an active game.
//
// The model implements the tea.Model interface and can be used with
// Bubble Tea for TUI rendering and event handling. The game state is
// immutable - updates create new GameModel instances.
//
// Fields:
//   - Config: Game configuration (board size, speed, level)
//   - Snake: Positions of snake body segments (head is first)
//   - Food: Current food location and type
//   - Direction: Current movement direction
//   - Score: Current game score
//   - IsGameOver: Whether the game has ended
//   - IsOutOfBounds: Whether snake hit a boundary
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

#### For Interfaces

```go
// ScoreService defines operations for managing game scores.
//
// Implementations are responsible for persistent storage and retrieval of
// score records. All methods should support context cancellation via ctx.
type ScoreService interface {
    // GetHighScore returns the highest score ever recorded
    GetHighScore(ctx context.Context) (Score, error)
    
    // GetScores returns all recorded scores, ordered by value descending
    GetScores(ctx context.Context) ([]Score, error)
    
    // SetCurrentScore saves the final score for the current session
    SetCurrentScore(ctx context.Context, value int) error
    
    // GetCurrentScore returns the score accumulated so far in current session
    GetCurrentScore(ctx context.Context) (int, error)
}
```

#### For Constants/Enums

```go
// Direction represents movement direction for the snake.
type Direction int

const (
    // Up means move toward top of screen
    Up Direction = iota
    // Down means move toward bottom of screen
    Down
    // Left means move toward left of screen
    Left
    // Right means move toward right of screen
    Right
)
```

---

## Complex Logic Documentation

For complex algorithms or non-obvious logic, add inline comments explaining the logic.

### Pattern 1: Algorithm Overview

```go
func (g *GameModel) moveSnake() {
    // Algorithm: Remove tail, add new head position
    // This maintains a fixed snake length while moving it forward
    // and allows for growth when food is consumed.
    
    // Calculate next position based on current direction
    nextPos := g.calculateNextPosition()
    
    // Add new position at head of snake
    g.Snake = append([]Position{nextPos}, g.Snake...)
    
    // Remove tail segment (unless we just ate food)
    if !g.justAteFood {
        g.Snake = g.Snake[:len(g.Snake)-1]
    }
}
```

### Pattern 2: Complex Condition

```go
func (g *GameModel) checkCollision() bool {
    headPos := g.Snake[0]
    
    // Check bounds: if head is outside game area
    if headPos.X < 0 || headPos.X >= g.Config.Columns ||
       headPos.Y < 0 || headPos.Y >= g.Config.Rows {
        return true
    }
    
    // Check self-collision: if head overlaps body (skip index 0 which is head)
    for i := 1; i < len(g.Snake); i++ {
        if g.Snake[i] == headPos {
            return true
        }
    }
    
    // Check obstacle collision: if head hits a pillar
    for _, pillar := range g.Config.Pillars {
        if pillar == headPos {
            return true
        }
    }
    
    return false
}
```

### Pattern 3: Non-obvious Loop Logic

```go
// Regenerate food at random location until we find an empty spot
func (g *GameModel) spawnFood() {
    for {
        newFood := Food{
            Position: Position{
                X: rand.Intn(g.Config.Columns),
                Y: rand.Intn(g.Config.Rows),
            },
        }
        
        // Ensure food doesn't spawn on snake body or obstacle
        if !g.isSnake(newFood.X, newFood.Y) && !g.isObstacle(newFood.X, newFood.Y) {
            g.Food = newFood
            break
        }
    }
}
```

---

## Examples

Real examples from the codebase with proper documentation:

### Example 1: Well-Documented Type

```go
// Score represents a single game score entry in the database.
//
// Score records are immutable once created and include metadata about
// when and by whom the score was achieved. Scores are ordered by Value
// in leaderboards.
//
// Fields:
//   - ID: Unique identifier assigned by database
//   - User: Hostname of player who achieved the score
//   - Session: Unique ID for the game session in which score was achieved
//   - Value: The numerical score (typically food count)
//   - CreatedAt: Timestamp of when the score was recorded
type Score struct {
    ID        int       `db:"id"`
    User      string    `db:"user"`
    Session   string    `db:"session"`
    Value     int       `db:"value"`
    CreatedAt time.Time `db:"created_at"`
}
```

### Example 2: Well-Documented Interface

```go
// SessionManager handles creation and tracking of game sessions.
//
// Each game session gets a unique identifier generated from cryptographically
// random bytes. Sessions are used to associate scores with individual games
// and prevent duplicate score entries.
type SessionManager interface {
    // CreateNewSession creates a new unique session identifier.
    // The value parameter is reserved for future use.
    // Returns a 20-character base64-encoded random string.
    CreateNewSession(value any) (string, error)
    
    // DestroyCurrentSession clears the active session.
    // Typically called when a game ends.
    DestroyCurrentSession() error
    
    // GetCurrentSession returns the active session ID.
    // Returns an error if no session is currently active.
    GetCurrentSession() (string, error)
}
```

### Example 3: Well-Documented Method

```go
// Update processes input messages and updates the game state according to MVU pattern.
//
// For each message, the appropriate handler is called which may:
// - Update direction based on keyboard input
// - Move the snake on timer tick
// - Check for collisions and game end conditions
// - Generate new food if current food was eaten
//
// The returned Model is a new GameModel with updated state (not mutating the receiver).
// The returned Cmd may contain operations like screen clears or quit commands.
//
// Messages handled:
//   - tea.KeyMsg: Parse arrow keys/WASD to change direction
//   - GameTickMsg: Move snake one position
//   - tea.WindowSizeMsg: Update board dimensions
func (g *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input...
    case GameTickMsg:
        // Move snake and check collisions...
    }
    return g, nil
}
```

---

## Best Practices

### Do's ✅

1. **Document All Public Types and Functions**
   ```go
   // Good: Every exported identifier has documentation
   // BadFunction would cause a linting error
   func GoodFunction() string {
       return "documented"
   }
   ```

2. **Explain the "Why", Not Just the "What"**
   ```go
   // Good
   // We check self-collision separately from boundary collision because
   // self-collision has higher performance requirements (O(n) vs O(1))
   
   // Bad
   // Check for self collision
   ```

3. **Document Edge Cases and Invariants**
   ```go
   // Good
   // Note: Snake must always have at least one segment (the head).
   // The slice is always maintained with head at index 0.
   type GameModel struct {
       Snake []Position
   }
   
   // Bad: No documentation of constraints
   ```

4. **Include Examples for Complex Functions**
   ```go
   // Good
   // Example:
   //   config := game.Level3GameConfig()
   //   model := game.InitalGameModel(config)
   ```

5. **Use Proper Grammar**
   ```go
   // Good
   // Score represents a game score record.
   
   // Bad
   // score - A game score record
   ```

### Don'ts ❌

1. **Don't Repeat Code as Documentation**
   ```go
   // Bad: Comment just repeats code
   i := 0  // set i to 0
   
   // Good: Explain why, not what
   i := 0  // Use 0-based index for array access
   ```

2. **Don't Let Comments Get Out of Sync**
   ```go
   // Bad: This gets outdated
   // Returns the high score (returns all scores now)
   func (s *ScoreServiceImpl) GetHighScore(ctx context.Context) ([]Score, error) {
   
   // Good: Update comments when code changes
   ```

3. **Don't Over-Document Trivial Code**
   ```go
   // Bad: Too detailed for simple getter
   // This function returns the current score
   // The score is the current value
   // It comes from the Score field
   func (g *GameModel) GetScore() int {
       return g.Score
   }
   
   // Good: Keep it brief
   // GetScore returns the current game score.
   func (g *GameModel) GetScore() int {
   ```

4. **Don't Use Abbreviations**
   ```go
   // Bad: Hard to understand
   // VerifyPos checks if x, y is valid
   
   // Good: Clear language
   // VerifyPosition checks if x, y is within game bounds
   ```

---

## Documentation Coverage Checklist

For every public (capitalized) identifier, verify:

- [ ] Package has documentation comment
- [ ] Type has documentation comment
- [ ] All exported fields have descriptions
- [ ] All exported methods have documentation comments
- [ ] All exported functions have documentation comments
- [ ] Documentation explains the "why" not just the "what"
- [ ] Examples provided for complex functions
- [ ] Edge cases and constraints documented
- [ ] Documentation is grammatically correct

---

## Running Documentation Checks

### View Generated Documentation

```bash
# View package documentation
go doc github.com/the-Jinxist/golang_snake_game/internal

# View function documentation
go doc github.com/the-Jinxist/golang_snake_game/internal.GetScoreService

# Generate HTML docs
godoc -http=:6060
# Then visit http://localhost:6060
```

### Check for Missing Documentation

```bash
# Using golangci-lint (if installed)
golangci-lint run ./... --disable-all --enable=revive

# Or using go vet
go vet ./...
```

---

## Template for New Functions

```go
// FunctionName brief description of what it does.
//
// Longer explanation if needed. This can describe the algorithm,
// when it should be used, or any important constraints.
//
// Parameters:
//   - param1: What this parameter does
//   - param2: What this parameter does
//
// Returns:
//   - ReturnType: What this return value represents
//   - error: Description of error conditions (if applicable)
//
// Example:
//   result, err := FunctionName(arg1, arg2)
//   if err != nil {
//       // handle error
//   }
func FunctionName(param1 string, param2 int) (string, error) {
    // implementation here
    return "", nil
}
```

---

**Remember**: Good documentation is written for the reader, not the writer. Assume your documentation will be read by someone unfamiliar with the code, perhaps months or years later.

Good documentation = Better code quality = Easier maintenance + Faster onboarding 📚
