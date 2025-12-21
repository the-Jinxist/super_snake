package game

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-Jinxist/golang_snake_game/internal"
)

type Tick struct{}

type TriggerNextLevel struct{}
type GameStartConfig struct {
	Rows           int
	Columns        int
	Pillars        []Position
	IsWalled       bool
	Level          int
	IsFinalLevel   bool
	ScoreThreshold int
	Scoring        int
	ScoreService   internal.ScoreService
	SessionManager internal.SessionManager
}

func TickGame() tea.Cmd {
	return func() tea.Msg {
		return Tick{}
	}
}

func DefaultGameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           30,
		Columns:        25,
		Scoring:        10,
		IsWalled:       true,
		ScoreThreshold: 200,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level1GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           35,
		Columns:        25,
		ScoreThreshold: 700,
		Scoring:        10,
		Pillars:        level1Pillars,
		Level:          1,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level2GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           40,
		Columns:        30,
		IsWalled:       true,
		ScoreThreshold: 1900,
		Scoring:        10,
		Level:          2,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level3GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           45,
		Columns:        35,
		ScoreThreshold: 3500,
		Scoring:        10,
		Level:          3,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level4GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           50,
		Columns:        40,
		ScoreThreshold: 5500,
		IsWalled:       true,
		Scoring:        10,
		Level:          4,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level5GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           60,
		Columns:        40,
		ScoreThreshold: 8000,
		Scoring:        10,
		IsWalled:       true,
		Level:          5,
		IsFinalLevel:   true,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}
