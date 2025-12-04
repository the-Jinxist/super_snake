package game

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &GameModel{}

type Position struct {
	X int
	Y int
}
type GameModel struct {
	Rows    int
	Columns int
	Snake   [][]int

	Food     Position
	Velocity Position
}

func InitalGameModel() *GameModel {
	return &GameModel{
		Rows:    50,
		Columns: 30,
		Velocity: Position{
			X: 1,
			Y: 0,
		},
	}
}

// Init implements tea.Model.
func (g *GameModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (g *GameModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return g, nil

}

// View implements tea.Model.
func (g *GameModel) View() string {

	cellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	var output string
	for i := 0; i < g.Columns; i++ {
		for j := 0; j < g.Rows; j++ {
			output += cellStyle.Render(EmptyCell)
		}
		output += "\n"
	}

	output = lipgloss.NewStyle().Margin(5).Border(lipgloss.DoubleBorder(), true).Render(output)
	return output
}
