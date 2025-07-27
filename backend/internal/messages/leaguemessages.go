package messages

type CreateLeagueRequest struct {
	LeagueName string `json:"league_name"`
	URL        string `json: "url"`
}

type ListLeagueResponseItem struct {
	LeagueName string `json:"league_name"`
	URL        string `json:"url"`
}

type ListLeagueResponse struct {
	Leagues []ListLeagueResponseItem `json: "leagues"`
}
