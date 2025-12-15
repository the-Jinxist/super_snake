package internal

import (
	"context"
	"database/sql"
	"log"
	"time"
)

var _ ScoreService = &ScoreServiceImpol{}

type Score struct {
	User      string    `db:"user"`
	Session   string    `db:"session"`
	Value     int       `db:"value"`
	CreatedAt time.Time `db:"created_at"`
}

type ScoreService interface {
	GetHighScore(ctx context.Context) (Score, error)
	GetScores(ctx context.Context) ([]Score, error)
	SetScore(ctx context.Context, value int) error
}

func NewScoreService(user, session string, db *sql.DB) ScoreService {
	return &ScoreServiceImpol{
		db:          db,
		CurrentUser: user,
		Session:     session,
	}
}

type ScoreServiceImpol struct {
	CurrentUser string
	db          *sql.DB
	Session     string
}

// GetHighScore implements ScoreService.
func (s *ScoreServiceImpol) GetHighScore(ctx context.Context) (Score, error) {
	var score Score

	rows, err := s.db.QueryContext(ctx, `select 1 from "scores" order by value`)
	if err != nil {
		return score, nil
	}

	err = rows.Scan(&score)
	if err != nil {
		return score, err
	}

	return score, nil
}

// GetScores implements ScoreService.
func (s *ScoreServiceImpol) GetScores(ctx context.Context) ([]Score, error) {

	var scores []Score

	rows, err := s.db.QueryContext(ctx, `select * from scores order by value limit 1`)
	if err != nil {
		return scores, nil
	}

	defer rows.Close()
	for rows.Next() {
		var score Score
		err = rows.Scan(&score)
		if err != nil {
			log.Printf("issue with scanning each rows: %s", err)
			return scores, nil
		}

		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error encountered while iterating over rows: %s", err)
	}

	return scores, nil

}

// SetScore implements ScoreService.
func (s *ScoreServiceImpol) SetScore(ctx context.Context, value int) error {
	_, err := s.db.ExecContext(ctx,
		`insert into scores ("user", session, value) 
	 values (?, ?, ?)
	 on conflict(session) do update set
		value = excluded.value
	where scores.session = excluded.session and scores."user" = excluded."user";
	 `, s.CurrentUser, s.Session, value)
	if err != nil {
		return err
	}

	return nil
}
