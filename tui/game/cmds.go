package game

import (
	"time"

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
	IsDebugGrid    bool
	FPS            time.Duration
	ScoreService   internal.ScoreService
	SessionManager internal.SessionManager
}

func TickGame() tea.Cmd {
	return func() tea.Msg {
		return Tick{}
	}
}

func DebugGameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           30,
		Columns:        25,
		Scoring:        10,
		IsWalled:       true,
		IsDebugGrid:    true,
		ScoreThreshold: 200,
		FPS:            time.Millisecond * 250,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func DefaultGameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           30,
		Columns:        25,
		Scoring:        10,
		IsWalled:       true,
		FPS:            time.Millisecond * 250,
		ScoreThreshold: 20, //TODO MUST REMOVE
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
		FPS:            time.Millisecond * 200,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level2GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           35,
		Columns:        25,
		IsWalled:       true,
		ScoreThreshold: 1900,
		Scoring:        10,
		Level:          2,
		FPS:            time.Millisecond * 200,
		Pillars:        level2Pillars,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level3GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           35,
		Columns:        25,
		ScoreThreshold: 3500,
		Scoring:        10,
		Level:          3,
		FPS:            time.Millisecond * 150,
		Pillars:        level3Pillars,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level4GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           35,
		Columns:        25,
		ScoreThreshold: 5500,
		IsWalled:       true,
		Pillars:        level1Pillars,
		Scoring:        10,
		Level:          4,
		FPS:            time.Millisecond * 150,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}

func Level5GameConfig() GameStartConfig {
	return GameStartConfig{
		Rows:           35,
		Columns:        25,
		ScoreThreshold: 8000,
		Scoring:        8000,
		IsWalled:       true,
		Level:          5,
		IsFinalLevel:   true,
		FPS:            time.Millisecond * 150,
		Pillars:        level3Pillars,
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}
