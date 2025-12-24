package game

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/Broderick-Westrope/charmutils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/the-Jinxist/golang_snake_game/tui/views"
	"github.com/the-Jinxist/golang_snake_game/utils"
)

var _ tea.Model = &GameModel{}

type Direction int

type Position struct {
	X int
	Y int
}

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Food struct {
	Position
	BigFish bool
}

type GameModel struct {
	Config GameStartConfig
	Snake  []Position

	Food          Food
	Direction     Direction
	Score         int
	IsGameOver    bool
	IsOutOfBounds bool
	spinner       spinner.Model
	rng           *rand.Rand
	isPaused      bool
}

func InitalGameModel(gameConfig GameStartConfig) *GameModel {
	source := rand.NewPCG(uint64(time.Now().Unix()), uint64(time.Now().UnixMicro()))

	s := spinner.New()
	s.Spinner = spinner.Dot

	currentScore, err := gameConfig.ScoreService.GetCurrentScore(context.Background())
	if err != nil {
		log.Fatalf("Failed to get current score: %s", err)
	}

	gameMod := &GameModel{
		Snake: []Position{
			{
				X: gameConfig.Rows / 2,
				Y: gameConfig.Columns / 2,
			},
		},
		rng:       rand.New(source),
		Config:    gameConfig,
		Direction: Right,
		Score:     currentScore,
		spinner:   s,
	}

	return gameMod
}

func (g *GameModel) isSnake(x int, y int) bool {
	isSnake := false
	for _, snakeBody := range g.Snake {
		if snakeBody.X == x && snakeBody.Y == y {
			isSnake = true
		}
	}

	return isSnake
}

func (g *GameModel) isSnakeHead(x int, y int) bool {
	snakeHead := g.Snake[0]

	return x == snakeHead.X && y == snakeHead.Y
}

func (g *GameModel) isFood(x int, y int) bool {
	food := g.Food
	return x == food.X && y == food.Y
}

func (g *GameModel) directionToPosition(direction Direction) Position {
	position := Position{
		Y: 0, X: 1,
	}

	switch direction {
	case Up:
		position = Position{Y: -1, X: 0}
	case Down:
		position = Position{Y: 1, X: 0}
	case Left:
		position = Position{Y: 0, X: -1}
	}

	return position
}

func (g *GameModel) hasHitWall(x int, y int) bool {
	return g.isOutOfBounds(x, y) && g.Config.IsWalled
}

func (g *GameModel) isPillar(x, y int) bool {
	for _, pos := range g.Config.Pillars {
		if x == pos.X && y == pos.Y {
			return true
		}
	}

	return false
}

func (g *GameModel) isOutOfBounds(x int, y int) bool {
	return x > g.Config.Rows-1 || y > g.Config.Columns-1 || x < 0 || y < 0
}

// Init implements tea.Model.
func (g *GameModel) Init() tea.Cmd {
	g.instantiateFood()
	return tea.Batch(g.Tick())
}

func (g *GameModel) instantiateFood() {

	rows := g.Config.Rows
	columns := g.Config.Columns

	randomX := g.rng.IntN(rows)
	randomY := g.rng.IntN(columns)

	for g.isSnake(randomX, randomY) || g.isPillar(randomX, randomY) || g.hasHitWall(randomX, randomY) {
		randomX = g.rng.IntN(rows)
		randomY = g.rng.IntN(columns)
	}

	g.Food = Food{
		Position: Position{
			X: randomX,
			Y: randomY,
		},
	}
}

// Update implements tea.Model.
func (g *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if g.Config.IsDebugGrid {
		return g, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		input := msg.String()

		if g.hasReachedLevelThreshold() {
			return g, nil
		}

		if g.IsGameOver {
			g.Config.SessionManager.DestroyCurrentSession()
			if utils.KeyMatchesInput(input, utils.Esc, utils.Space) {
				return g, tea.Batch(views.ClearScreen(), views.SwitchModeCmd(views.ModeMenu))
			}

			return g, nil
		}

		if utils.KeyMatchesInput(input, utils.Space) {
			g.isPaused = !g.isPaused
		}

		if g.isPaused {

			if utils.KeyMatchesInput(input, utils.Esc) {
				return g, tea.Batch(views.SwitchModeCmd(views.ModeMenu))
			}

			return g, nil
		}

		if utils.KeyMatchesInput(input, utils.KeyUp) {
			if g.Direction == Left || g.Direction == Right {
				g.Direction = Up
			}
		}

		if utils.KeyMatchesInput(input, utils.KeyRight) {
			if g.Direction == Up || g.Direction == Down {
				g.Direction = Right
			}
		}

		if utils.KeyMatchesInput(input, utils.KeyDown) {
			if g.Direction == Left || g.Direction == Right {
				g.Direction = Down
			}
		}

		if utils.KeyMatchesInput(input, utils.KeyLeft) {
			if g.Direction == Up || g.Direction == Down {
				g.Direction = Left
			}
		}

		return g, nil

	case Tick:

		if !g.isPaused && !g.IsGameOver {
			g.moveSnake()
		}

		return g, tea.Batch(g.Tick())
	default:
		return g, tea.Batch(g.Tick())

	}

}

func (g *GameModel) hasReachedLevelThreshold() bool {
	return g.Score == g.Config.ScoreThreshold
}

func (g *GameModel) Tick() tea.Cmd {

	if g.Config.IsDebugGrid {
		return nil
	}

	if g.hasReachedLevelThreshold() {

		if g.Config.Level == 5 {
			fmt.Print("\033[H\033[2J")
			return tea.Batch(views.SwitchModeCmd(views.ModeGameCompleted))
		}

		time.Sleep(2 * time.Second)
		nextLevel := views.NextLevelModeFromCurrent(g.Config.Level)

		fmt.Print("\033[H\033[2J")
		return tea.Batch(views.SwitchModeCmd(nextLevel))
	}

	return tea.Tick(g.Config.FPS, func(t time.Time) tea.Msg {
		return Tick{}
	})
}

func (g *GameModel) moveSnake() {

	if g.Config.IsDebugGrid {
		return
	}

	pos := g.directionToPosition(g.Direction)

	currentSnakeHead := g.Snake[0]

	movingToX := currentSnakeHead.X + pos.X
	movingToY := currentSnakeHead.Y + pos.Y

	//If snake is out of bounds
	if g.hasHitWall(movingToX, movingToY) {
		g.IsGameOver = true
	}

	if g.isPillar(movingToX, movingToY) {
		g.IsGameOver = true
	}

	// If snakes eats itself
	if g.isSnake(movingToX, movingToY) {
		g.IsGameOver = true
	}

	if g.IsGameOver || g.hasReachedLevelThreshold() {
		return
	}

	if g.isFood(movingToX, movingToY) {
		g.instantiateFood()
		newSnakeHead := Position{
			X: movingToX,
			Y: movingToY,
		}

		g.Snake = append([]Position{newSnakeHead}, g.Snake...)

		g.increaseScore()

	}

	//Snap snake back to the opp. side of stage if he goes out of bounds
	if g.isOutOfBounds(movingToX, movingToY) ||
		g.isOutOfBounds(currentSnakeHead.X, currentSnakeHead.Y) {

		// g.IsOutOfBounds = true
		if movingToX > g.Config.Rows-1 || currentSnakeHead.X > g.Config.Rows-1 {
			movingToX = 0
		}

		if movingToX < 0 || currentSnakeHead.X < 0 {
			movingToX = g.Config.Rows - 1
		}

		if movingToY > g.Config.Columns-1 || currentSnakeHead.Y > g.Config.Columns-1 {
			movingToY = 0
		}

		if movingToY < 0 || currentSnakeHead.Y < 0 {
			movingToY = g.Config.Columns - 1
		}
	}

	newSnakeHead := Position{
		X: movingToX,
		Y: movingToY,
	}

	g.Snake = g.Snake[:len(g.Snake)-1]
	g.Snake = append([]Position{newSnakeHead}, g.Snake...)
}

func (g *GameModel) increaseScore() {
	g.Score += g.Config.Scoring

	go func() {
		g.Config.ScoreService.SetCurrentScore(context.Background(), g.Score)
	}()

}

// View implements tea.Model.
func (g *GameModel) View() string {

	cellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	var output string
	for i := range g.Config.Columns {
		for j := range g.Config.Rows {

			if g.Config.IsDebugGrid {
				output += lipgloss.NewStyle().
					Height(5).
					Width(5).
					Border(lipgloss.BlockBorder(), true).
					Render(fmt.Sprintf("[%d,%d]", j, i))
				continue
			}

			if g.isSnake(j, i) {
				if g.isSnakeHead(j, i) {
					output += SnakeHeadFromDirection(g.Direction)
				} else {
					output += cellStyle.Render(FilledCell)
				}

			} else if g.isFood(j, i) {
				output += FoodCell()
			} else if g.isPillar(j, i) {
				output += PillarCell
			} else {
				output += cellStyle.Render(EmptyCell)
			}

		}
		output += "\n"
	}

	output = lipgloss.NewStyle().Border(lipgloss.ASCIIBorder(), g.Config.IsWalled).Render(output)

	output += "\n"

	if g.isPaused {
		output += lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Render(fmt.Sprintf("[ PAUSED ]. Your score: %d/%d. Press SPACE to resume! Press ESC to back to menu", g.Score, g.Config.ScoreThreshold))
	} else {
		output += lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Render(fmt.Sprintf("Your score: %d/%d. Press SPACE to pause!", g.Score, g.Config.ScoreThreshold))
	}

	if g.hasReachedLevelThreshold() {
		levelingUpMsg := "We're going up!"
		levelingUpMsg += "/n"
		levelingUpMsg += lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Render(g.spinner.View())
		output, _ = charmutils.OverlayCenter(output, levelingUpMsg, false)

	}

	if g.IsGameOver {
		gameOverMessage := gameOverMsg
		gameOverMessage += "\n"
		gameOverMessage += lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Render(fmt.Sprintf("Your final score is %d/%d\nPress SPACE to go back to menu", g.Score, g.Config.ScoreThreshold))
		output, _ = charmutils.OverlayCenter(output, gameOverMessage, false)
	}

	levelIndicator := generateLevelIndicator(g.Config.Level)

	help := generateHelpString()
	help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444745")).
		Render(help)

	return levelIndicator + output + "\n" + help
}

func generateHelpString() string {
	help := "\n[INSTRUCTIONS]:\n · -> or D to move right\n · <- or A to move left\n · ↑ or W to move up\n · ↓ or S to move down"

	if utils.IsWindowsMachine() {
		help = strings.ReplaceAll(help, "\n", " | ")
	}

	return help
}

func generateLevelIndicator(level int) string {
	lvlString := fmt.Sprintf("Level %d", level)
	if !utils.IsWindowsMachine() {
		lvlString = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(1).
			Background(lipgloss.Color("#3297a8")).
			Render(lvlString)
	}

	return lvlString + "\n"
}
