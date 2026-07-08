package models

type Vote struct {
	UserID     int `json:"user_id"`
	VotingID   int `json:"voting_id"`
	VoteOption int `json:"vote_option"`
}
