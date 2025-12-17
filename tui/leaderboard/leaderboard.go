package leaderboard

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/the-Jinxist/golang_snake_game/internal"
	"github.com/the-Jinxist/golang_snake_game/tui/views"
)

const listHeight = 14

var (
	_ tea.Model = &Leaderboard{}

	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type item struct {
	TimeOfScore time.Time
	Score       int
}

func (i item) FilterValue() string { return "" }
func (i item) Title() string       { return fmt.Sprintf("%d", i.Score) }
func (i item) Description() string {
	return fmt.Sprintf("Score was recorded at: %s", i.TimeOfScore)
}

type Leaderboard struct {
	list   list.Model
	Scores []internal.Score
	Config LeaderboardConfig
}

func NewLeaderboardModel(config LeaderboardConfig) *Leaderboard {

	const defaultWidth = 20
	var userItems []list.Item
	scores, _ := config.ScoreService.GetScores(context.Background())

	for _, i := range scores {
		userItems = append(userItems, item{
			TimeOfScore: i.CreatedAt,
			Score:       i.Value,
		})
	}

	l := list.New(userItems, list.DefaultDelegate{}, defaultWidth, listHeight)
	l.Title = "Your leaderboard!"
	l.Styles.Title = titleStyle

	return &Leaderboard{
		Config: config,
	}
}

// Init implements tea.Model.
func (l *Leaderboard) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (l *Leaderboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.list.SetWidth(msg.Width)
		return l, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return l, tea.Quit
		case "esc":
			return l, tea.Batch(views.SwitchModeCmd(views.ModeMenu))
		}
	}

	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

// View implements tea.Model.
func (l *Leaderboard) View() string {
	return "\n" + l.list.View()
}
