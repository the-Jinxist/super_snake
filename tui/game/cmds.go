package game

import tea "github.com/charmbracelet/bubbletea"

type Tick struct{}
type GameStartConfig struct {
	Rows    int
	Columns int
}

func TickGame() tea.Cmd {
	return func() tea.Msg {
		return Tick{}
	}
}

func DefaultGameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:    40,
		Columns: 25,
	}
}
