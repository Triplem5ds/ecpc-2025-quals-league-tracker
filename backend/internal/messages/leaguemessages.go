package messages

type LeagueScoring struct {
	Rank  int `json:"rank"`
	Score int `json:"score"`
}

type LeaguePlayers struct {
	IdLeaguePlayer int    `json:"id_league_player"`
	PlayerName     string `json:"player_name"`
}

type LeagueRounds struct {
	IdLeagueRound int    `json:"id_league_round"`
	RoundName     string `json:"round_name"`
}

type LeagueScoreboard struct {
	Rank         int    `json:"rank"`
	PlayerName   string `json:"player_name"`
	RoundsPlayed int    `json:"rounds_played"`
	TotalPoints  int    `json:"total_points"`
}

type LeagueRoundsResults struct {
	IdLeaguePlayer int `json:"id_league_player"`
	IdLeagueRound  int `json:"id_league_round"`
	Rank           int `json:"rank"`
}

type CreateLeagueRequest struct {
	LeagueName  string          `json:"league_name"`
	Description string          `json:"description"`
	Scoring     []LeagueScoring `json:"scoring"`
}

type ListLeagueResponseItem struct {
	IdLeague    int    `json:"id_league"`
	LeagueName  string `json:"league_name"`
	Description string `json:"url"`
}

type ListLeagueResponse struct {
	Leagues []ListLeagueResponseItem `json:"leagues"`
}

type GetLeagueItem struct {
	IdLeague    int             `json:"id_league"`
	LeagueName  string          `json:"league_name"`
	Description string          `json:"description"`
	Scoring     []LeagueScoring `json:"scoring"`
	Rounds      []LeagueRounds  `json:"rounds"`
	Players     []LeaguePlayers `json:"players"`
}

type GetLeagueResponse struct {
	League GetLeagueItem `json:"league"`
}

type GetLeagueScoreboardResponse struct {
	Scoreboard    []LeagueScoreboard    `json:"scoreboard"`
	Rounds        []LeagueRounds        `json:"league_rounds"`
	Players       []LeaguePlayers       `json:"league_players"`
	RoundsResults []LeagueRoundsResults `json:"rounds_results"`
}
