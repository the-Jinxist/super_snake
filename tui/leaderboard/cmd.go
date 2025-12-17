package leaderboard

import "github.com/the-Jinxist/golang_snake_game/internal"

type LeaderboardConfig struct {
	ScoreService   internal.ScoreService
	SessionManager internal.SessionManager
}

func DefaultLeaderboardConfig() LeaderboardConfig {
	return LeaderboardConfig{
		ScoreService:   internal.GetScoreService(),
		SessionManager: internal.GetSessionManager(),
	}
}
