package messages

import "time"

type RoundScoring struct {
	IdLeaguePlayer int    `json:"id_league_player"`
	PlayerName     string `json:"player_name"`
	Rank           int    `json:"rank"`
}

type CreateRoundRequest struct {
	IdLeague    int       `json:"id_league"`
	RoundName   string    `json:"round_name"`
	URL         string    `json:"round_url"`
	StartTime   time.Time `json:"start_time"`
	Description string    `json:"description"`
}

type UpdateRoundRequest struct {
	IdLeagueRound int            `json:"id_league_round"`
	Scoring       []RoundScoring `json:"scoring"`
}

type GetRoundItem struct {
	IdLeagueRound int            `json:"id_league_round"`
	RoundName     string         `json:"round_name"`
	URL           string         `json:"round_url"`
	StartTime     time.Time      `json:"start_time"`
	Description   string         `json:"description"`
	Scoring       []RoundScoring `json:"scoring"`
}

type GetRoundResponse struct {
	Round GetRoundItem `json:"round"`
}
