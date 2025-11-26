package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
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
	`
	_, err := d.conn.Exec(query, userID, votingID, voteOption, time.Now())
	return err
}

func (d *DB) Close() error {
	return d.conn.Close()
}
