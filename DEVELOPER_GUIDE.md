# Developer Guide & Contributing

This guide provides information for developers who want to understand, modify, or contribute to the Super Snake project.

---

## Table of Contents

- [Development Setup](#development-setup)
- [Project Structure Overview](#project-structure-overview)
- [Building & Running](#building--running)
- [Code Organization](#code-organization)
- [Common Development Tasks](#common-development-tasks)
- [Testing Strategy](#testing-strategy)
- [Git Workflow](#git-workflow)
- [Coding Standards](#coding-standards)
- [Debugging Tips](#debugging-tips)
- [Performance Tips](#performance-tips)
- [Contributing Guidelines](#contributing-guidelines)

---

## Development Setup

### Prerequisites

- **Go**: 1.24.2 or higher
  ```bash
  go version  # Verify installation
  ```

- **Git**: For version control
  ```bash
  git version
  ```

- **Terminal**: With ANSI color support
  - macOS: Terminal, iTerm2, Alacritty
  - Linux: Most terminals support this
  - Windows: Windows Terminal (recommended)

### Environment Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/the-Jinxist/golang_snake_game.git
   cd golang_snake_game
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   go mod verify
   ```

3. **Verify Build**
   ```bash
   go build -o super_snake
   ```

4. **Set Up IDE (VS Code)**
   - Install "Go" extension (golang.Go)
   - Install "Bubble Tea" helpers (optional)
   - Set GOPROXY if behind corporate proxy

### Database Setup

The application uses SQLite with automatic initialization:

```go
// Database is created automatically on first run
// File: my.db (in project root)
// Schema: Created by internal/db.go during Init()
```

---

## Project Structure Overview

```
golang_snake_game/
├── main.go                          # Entry point
├── go.mod, go.sum                   # Dependencies
├── README.md                         # User documentation
├── ARCHITECTURE.md                   # Architecture documentation
├── API_REFERENCE.md                  # API reference
├── DEVELOPER_GUIDE.md               # This file
│
├── cmd/
│   └── root.go                      # Cobra CLI command
│
├── internal/
│   ├── internal.go                  # Service initialization
│   ├── db.go                        # Database setup
│   ├── session.go                   # Session management
│   └── score.go                     # Score service
│
├── tui/
│   ├── super_snake.go               # Root TUI model
│   ├── menu/
│   │   └── start_game.go            # Menu implementation
│   ├── game/
│   │   ├── game.go                  # Game logic
│   │   ├── game_over.go             # Game over screen
│   │   ├── cmds.go                  # Game messages
│   │   ├── styles.go                # UI styling
│   │   ├── pillars.go               # Level obstacles
│   │   └── pillar_two.go            # More obstacles
│   ├── leaderboard/
│   │   ├── leaderboard.go           # Leaderboard display
│   │   └── cmd.go                   # Leaderboard commands
│   └── views/
│       └── mode.go                  # View mode definitions
│
└── utils/
    └── keys.go                      # Keyboard utilities
```

---

## Building & Running

### Development Build

```bash
# Build executable
go build -o super_snake

# Run the game
./super_snake

# Or run directly
go run main.go
```

### Production Build

```bash
# Build optimized binary
go build -ldflags="-s -w" -o super_snake

# Build for specific OS
GOOS=linux GOARCH=amd64 go build -o super_snake-linux
GOOS=darwin GOARCH=amd64 go build -o super_snake-macos
GOOS=windows GOARCH=amd64 go build -o super_snake.exe
```

### Running with Debug Output

```bash
# Run with log output visible
go run main.go 2>&1 | tee debug.log

# Verbose error messages
GODEBUG=all go run main.go
```

---

## Code Organization

### Package Responsibilities

**cmd/**
- CLI application setup
- Cobra command configuration
- Service initialization trigger
- Application entry point coordination

**internal/**
- Service implementations
- Database management
- Business logic abstraction
- Global state management

**tui/**
- User interface implementation
- Bubble Tea models
- Screen rendering
- User input handling
- Navigation between views

**utils/**
- Utility functions
- Keyboard input handling
- Helper functions (non-package-specific)

### Dependency Layers

```
cmd (depends on) → internal, tui
tui (depends on) → internal, utils, views
views (depends on) → nothing (pure definitions)
internal (depends on) → database/sql, crypto
utils (depends on) → nothing (pure utilities)
```

**Rule**: Lower layers should not depend on higher layers (no circular dependencies).

---

## Common Development Tasks

### Add a New Game Level

1. **Create Configuration** (`tui/game/game.go`)
   ```go
   func Level6GameConfig() GameStartConfig {
       return GameStartConfig{
           ScoreService:   internal.GetScoreService(),
           SessionManager: internal.GetSessionManager(),
           Rows:           30,
           Columns:        60,
           SpeedMs:        50,  // Faster than Level 5
           Level:          5,
           Pillars:        level6Pillars,
       }
   }
   ```

2. **Define Obstacles** (`tui/game/pillars.go`)
   ```go
   var level6Pillars = []Position{
       // Define your obstacle positions
       {X: 10, Y: 10},
       {X: 11, Y: 10},
       // ... more positions
   }
   ```

3. **Add Mode** (`tui/views/mode.go`)
   ```go
   const (
       // ... existing modes
       ModeGame6 = Mode(iota)
   )
   
   func NextLevelModeFromCurrent(level int) Mode {
       if level == 5 {
           return ModeGame6
       }
       // ... rest of function
   }
   ```

4. **Route Configuration** (`tui/super_snake.go`)
   ```go
   func NextLevelConfigFromMode(level views.Mode) game.GameStartConfig {
       switch level {
       // ... existing cases
       case views.ModeGame6:
           return game.Level6GameConfig()
       }
   }
   ```

### Add a New Service

1. **Define Interface** (`internal/new_service.go`)
   ```go
   type NewService interface {
       DoSomething(ctx context.Context) error
       GetData(ctx context.Context) (string, error)
   }
   ```

2. **Implement Service**
   ```go
   type NewServiceImpl struct {
       db *sql.DB
       // ... other fields
   }
   
   func (s *NewServiceImpl) DoSomething(ctx context.Context) error {
       // implementation
       return nil
   }
   ```

3. **Register in Init** (`internal/internal.go`)
   ```go
   var (
       sessionManager SessionManager
       scoreService   ScoreService
       newService     NewService  // Add this
   )
   
   func GetNewService() NewService {
       return newService
   }
   
   func IntializeConfigs() {
       // ... existing initialization
       newService = NewNewService(db)  // Initialize it
   }
   ```

4. **Inject into Models**
   ```go
   type GameStartConfig struct {
       ScoreService   ScoreService
       SessionManager SessionManager
       NewService     NewService  // Add this
       // ...
   }
   ```

### Add a New View/Screen

1. **Create Package** (`tui/newview/`)
   ```go
   package newview
   
   import tea "github.com/charmbracelet/bubbletea"
   
   type NewViewModel struct {
       // state fields
   }
   
   func NewModel(config Config) *NewViewModel {
       return &NewViewModel{
           // initialize
       }
   }
   
   func (m *NewViewModel) Init() tea.Cmd {
       return nil
   }
   
   func (m *NewViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
       switch msg := msg.(type) {
       case tea.KeyMsg:
           // handle input
       }
       return m, nil
   }
   
   func (m *NewViewModel) View() string {
       return "Your view here"
   }
   ```

2. **Add Mode** (`tui/views/mode.go`)
   ```go
   const (
       // ... existing
       ModeNewView = Mode(iota)
   )
   ```

3. **Route in SuperSnake** (`tui/super_snake.go`)
   ```go
   func (s *SuperSnake) setChild(mode views.Mode) {
       switch mode {
       case views.ModeNewView:
           s.child = newview.NewModel(config)
           return
       // ... other cases
       }
   }
   ```

### Modify Game Rules

The game loop is in `tui/game/game.go`:

```go
func (g *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // This is where game mechanics happen:
    // 1. Process input
    // 2. Move snake
    // 3. Check collisions
    // 4. Spawn food
    // 5. Update score
}
```

Key methods to understand:
- `moveSnake()`: Update snake position
- `checkCollision()`: Detect collisions
- `spawnFood()`: Generate food location
- `isSnake()`: Check if position is snake body

---

## Testing Strategy

### Current Testing Coverage

The project currently lacks comprehensive tests. Here's how to add them:

### Unit Tests

**Example: Test keyboard input mapping** (`utils/keys_test.go`)

```go
package utils

import "testing"

func TestKeyMatchesInput(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        keys     []Key
        expected bool
    }{
        {"Match w to KeyUp", "w", []Key{KeyUp}, true},
        {"Match A to KeyUp", "A", []Key{KeyUp}, true},
        {"No match", "x", []Key{KeyUp}, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := KeyMatchesInput(tt.input, tt.keys...)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

**Run tests**:
```bash
go test ./...                    # All tests
go test -v ./utils              # Verbose utils tests
go test -cover ./...            # Coverage report
go test -coverprofile=cov.out ./... && go tool cover -html=cov.out
```

### Integration Tests

**Example: Test score service** (`internal/score_test.go`)

```go
package internal

import (
    "context"
    "testing"
    _ "github.com/mattn/go-sqlite3"
)

func TestScoreService(t *testing.T) {
    // Setup
    db := CreateDB()
    defer db.Close()
    
    sessionMgr := NewSessionManager()
    scoreService := NewScoreService("test_user", sessionMgr, db)
    
    // Test
    err := scoreService.SetCurrentScore(context.Background(), 100)
    if err != nil {
        t.Fatalf("SetCurrentScore failed: %v", err)
    }
    
    score, err := scoreService.GetHighScore(context.Background())
    if err != nil {
        t.Fatalf("GetHighScore failed: %v", err)
    }
    
    if score.Value != 100 {
        t.Errorf("got %d, want 100", score.Value)
    }
}
```

---

## Git Workflow

### Branch Naming

```
feature/[feature-name]      # New features
bugfix/[bug-description]    # Bug fixes
docs/[documentation-name]   # Documentation
refactor/[area]             # Code refactoring
test/[feature]              # Test additions
```

### Commit Messages

Follow conventional commits:

```
type(scope): description

[optional body]

[optional footer]
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

**Examples**:
```
feat(game): add level 6 with advanced obstacles
fix(score): prevent duplicate scores in database
docs(api): update API reference for ScoreService
refactor(tui): simplify game model initialization
test(utils): add comprehensive keyboard input tests
```

### Pull Request Process

1. Create feature branch from `main`
2. Make atomic, logical commits
3. Push to your fork
4. Create PR with clear description
5. Link related issues
6. Ensure CI passes
7. Request review from maintainers
8. Address feedback
9. Merge when approved

---

## Coding Standards

### Go Best Practices

1. **Naming Conventions**
   ```go
   // Good
   type GameModel struct { ... }
   func (g *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
   var currentScore int
   
   // Avoid
   type gamemodel struct { ... }
   type GameMdl struct { ... }
   var curr_score int
   var cs int
   ```

2. **Error Handling**
   ```go
   // Good
   result, err := operation()
   if err != nil {
       return fmt.Errorf("operation failed: %w", err)
   }
   
   // Avoid
   result, _ := operation()  // Ignoring errors
   if err != nil {
       panic(err)  // Use only for initialization
   }
   ```

3. **Interface Design**
   ```go
   // Good: Small, focused interfaces
   type Reader interface {
       Read(p []byte) (n int, err error)
   }
   
   // Avoid: Large, unfocused interfaces
   type Service interface {
       DoEverything() error
       AndMore() error
       AndMore2() error
   }
   ```

4. **Comments**
   ```go
   // Good: Clear, concise
   // Score represents a game score record
   type Score struct { ... }
   
   // Avoid: Redundant
   // Score struct
   type Score struct { ... }
   ```

### File Organization

- **One type per file** (or closely related types)
- **Tests in separate file** with `_test.go` suffix
- **Large functions/logic** extracted to separate files
- **Maximum 500 lines per file** (guideline)

### Import Organization

```go
import (
    // Standard library
    "context"
    "fmt"
    
    // External packages
    tea "github.com/charmbracelet/bubbletea"
    
    // Internal packages
    "github.com/the-Jinxist/golang_snake_game/internal"
)
```

---

## Debugging Tips

### Enable Debug Logging

```go
// In any file
import "log"

log.Printf("Debug: value=%v, state=%+v", value, model)
```

### Common Issues

**Problem**: Game not detecting collisions
```go
// Check Position struct calculations
// Ensure X and Y are not swapped
fmt.Printf("Snake head: %+v, Food: %+v\n", g.Snake[0], g.Food.Position)
```

**Problem**: Service returns nil
```go
// Ensure IntializeConfigs() was called
if internal.GetScoreService() == nil {
    log.Fatal("Services not initialized")
}
```

**Problem**: Database locked
```go
// Check if multiple connections to my.db
// SQLite doesn't handle concurrent writes well
```

### Using Go Debugger (Delve)

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug

# Set breakpoint
break main.main

# Continue execution
continue

# Step through
next
step
```

### Print-based Debugging

```go
// Temporary debug output
fmt.Fprintf(os.Stderr, "DEBUG: %+v\n", variable)

// In Bubble Tea (won't print to stdout)
log.SetOutput(os.Stderr)
log.Printf("Model state: %+v", m)
```

---

## Performance Tips

### Profiling

```bash
# CPU profile
go run -cpuprofile=cpu.prof main.go
go tool pprof cpu.prof

# Memory profile
go run -memprofile=mem.prof main.go
go tool pprof mem.prof

# Trace
go test -trace=trace.out ./...
go tool trace trace.out
```

### Optimization Opportunities

1. **Collision Detection** (currently O(n) for snake)
   ```go
   // Current: Iterate through all snake segments
   for _, segment := range g.Snake {
       if segment == position {
           return true
       }
   }
   
   // Better: Use map for O(1) lookup
   snakeMap := make(map[Position]bool)
   for _, segment := range g.Snake {
       snakeMap[segment] = true
   }
   if snakeMap[position] {
       return true
   }
   ```

2. **Database Queries** (no caching)
   ```go
   // Current: Query on every view
   scores, _ := scoreService.GetScores(ctx)
   
   // Better: Cache with expiration
   type CachedScoreService struct {
       delegate ScoreService
       cache    []Score
       expires  time.Time
   }
   ```

3. **Rendering** (rerender entire board)
   ```go
   // Current: Rebuild entire string every frame
   func (g *GameModel) View() string {
       return buildCompleteBoard(g)
   }
   
   // Better: Delta rendering (only changed cells)
   ```

---

## Contributing Guidelines

### Before You Start

1. Check [GitHub Issues](https://github.com/the-Jinxist/golang_snake_game/issues) for existing work
2. Discuss major changes in an issue first
3. Fork the repository
4. Create a feature branch

### Code Submission

1. **Style**: Run `go fmt ./...`
2. **Lint**: Use `golangci-lint` (if available)
3. **Tests**: Add tests for new functionality
4. **Documentation**: Update docs if behavior changes
5. **Commit messages**: Use conventional commits
6. **PR description**: Clear, detailed explanation

### What We're Looking For

✅ Bug fixes with test cases  
✅ New features with comprehensive documentation  
✅ Performance improvements with benchmarks  
✅ Code refactoring improving maintainability  
✅ Documentation improvements  

❌ Large, untested changes  
❌ Changes without clear motivation  
❌ Breaking changes without discussion  
❌ Code that doesn't follow Go conventions  

### Review Process

1. Automated checks (build, tests)
2. Code review by maintainers
3. Feedback discussion
4. Approval and merge

---

## Code Review Checklist

When reviewing code, check:

- [ ] Does it build without errors?
- [ ] Are tests included and passing?
- [ ] Does it follow Go conventions?
- [ ] Are comments clear and helpful?
- [ ] Is error handling appropriate?
- [ ] Are dependencies properly managed?
- [ ] Does it solve the stated problem?
- [ ] Could it cause performance issues?
- [ ] Is the API design sensible?
- [ ] Could it break existing functionality?

---

## Resources

### Go Resources

- [Go Documentation](https://golang.org/doc/)
- [Go Code Review Comments](https://golang.org/doc/effective_go)
- [GoLang Time Package](https://golang.org/pkg/time/)

### Bubble Tea

- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)

### Database

- [SQLite Go Bindings](https://github.com/mattn/go-sqlite3)
- [SQLite Documentation](https://www.sqlite.org/docs.html)

### Other

- [Cobra CLI Framework](https://cobra.dev/)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

## Questions?

- Open an [issue](https://github.com/the-Jinxist/golang_snake_game/issues)
- Check existing discussions
- Review the [API Reference](API_REFERENCE.md)
- Read the [Architecture](ARCHITECTURE.md) doc

---

**Happy coding! 🎮**
