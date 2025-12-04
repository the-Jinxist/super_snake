package game

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = GameModel{}

type GameModel struct{}

// Init implements tea.Model.
func (g GameModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (g GameModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return g, nil

}

// View implements tea.Model.
func (g GameModel) View() string {
	return "Game Mode!"
}
