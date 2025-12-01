package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/alex-carvalho/kafka-postgres-consumer/pkg/models"
)

type DB struct {
	conn *sql.DB
}

func Connect(url string) (*DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{conn: db}, nil
}

func (d *DB) InsertVote(userID, votingID, voteOption int) error {
	query := `
		INSERT INTO votes (user_id, voting_id, vote_option, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, voting_id) DO NOTHING
	`
	_, err := d.conn.Exec(query, userID, votingID, voteOption, time.Now())
	return err
}

func (d *DB) InsertVotesBatch(votes []models.Vote) error {
	if len(votes) == 0 {
		return nil
	}

	// Build the multi-value insert statement
	values := make([]string, len(votes))
	args := make([]interface{}, 0, len(votes)*4)

	for i, vote := range votes {
		values[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)",
			i*4+1, i*4+2, i*4+3, i*4+4)
		args = append(args, vote.UserID, vote.VotingID, vote.VoteOption, time.Now())
	}

	query := fmt.Sprintf(`
		INSERT INTO votes (user_id, voting_id, vote_option, created_at)
		VALUES %s
		ON CONFLICT (user_id, voting_id) DO NOTHING
	`, strings.Join(values, ","))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := d.conn.ExecContext(ctx, query, args...)
	return err
}

func (d *DB) Close() error {
	return d.conn.Close()
}
