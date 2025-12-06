package game

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type GameModel struct {
	Rows    int
	Columns int
	Snake   []Position

	Food      Position
	Direction Direction
	Ticker    *time.Ticker
}

func InitalGameModel() *GameModel {

	rows := 50
	column := 30

	gameMod := &GameModel{
		Rows:    rows,
		Columns: column,
		Snake: []Position{
			{
				X: rows / 2,
				Y: column / 2,
			},
		},
		Direction: Right,
		Ticker:    time.NewTicker(time.Microsecond * 6),
	}

	go manageTimer(gameMod.Ticker)
	return gameMod
}

func manageTimer(ticker *time.Ticker) {
	go func() {
		for {
			select {
			case _ = <-ticker.C:

			}
		}
	}()
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

func (g *GameModel) directionToPosition(direction Direction) Position {
	position := Position{
		Y: 0, X: 1,
	}

	switch direction {
	case Up:
		position = Position{Y: 1, X: 0}
	case Down:
		position = Position{Y: -1, X: 0}
	case Left:
		position = Position{Y: 0, X: -1}
	}

	return position
}

func (g *GameModel) isWall(x int, y int) bool {
	return x > g.Rows || y > g.Columns
}

// Init implements tea.Model.
func (g *GameModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (g *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd []tea.Cmd
	// fmt.Println("update!")
	// for range g.Ticker.C {
	g.moveSnake()
	// }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	case Tick:
		g.moveSnake()
	}

	return g, g.Tick()
}

func (g *GameModel) Tick() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return Tick{}
	})

}

func (g *GameModel) moveSnake() {

	pos := g.directionToPosition(g.Direction)

	currentSnakeHead := g.Snake[0]
	newSnakeHead := Position{
		X: currentSnakeHead.X + pos.X,
		Y: currentSnakeHead.Y + pos.Y,
	}

	g.Snake = g.Snake[:len(g.Snake)-1]
	g.Snake = append([]Position{newSnakeHead}, g.Snake...)
}

// View implements tea.Model.
func (g *GameModel) View() string {

	cellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	var output string
	for i := range g.Columns {
		for j := range g.Rows {

			if g.isSnake(j, i) {
				output += cellStyle.Render(FilledCell)
			} else {
				output += cellStyle.Render(EmptyCell)
			}

		}
		output += "\n"
	}

	output = lipgloss.NewStyle().Margin(5).Border(lipgloss.DoubleBorder(), true).Render(output)
	return output
}
