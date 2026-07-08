package models

type VoteRequest struct {
	UserID     int `json:"user_id" binding:"required,gt=0"`
	VotingID   int `json:"voting_id" binding:"required,gt=0"`
	VoteOption int `json:"vote_option" binding:"required,gte=0"`
}

type VoteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResultsResponse struct {
	VotingID   int         `json:"voting_id"`
	Results    map[int]int `json:"results"`
	TotalVotes int         `json:"total_votes"`
}
