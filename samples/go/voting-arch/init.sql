CREATE TABLE IF NOT EXISTS votes (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  voting_id INTEGER NOT NULL,
  vote_option INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, voting_id)
);

CREATE INDEX idx_voting_id ON votes(voting_id);
CREATE INDEX idx_user_voting ON votes(user_id, voting_id);
