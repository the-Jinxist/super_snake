package game

import tea "github.com/charmbracelet/bubbletea"

type Tick struct{}

func TickGame() tea.Cmd {
	return func() tea.Msg {
		return Tick{}
	}
}
