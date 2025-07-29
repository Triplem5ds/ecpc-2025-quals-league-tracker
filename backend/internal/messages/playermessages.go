package messages

type CreatePlayerRequest struct {
	IdLeague   int    `json:"id_league"`
	PlayerName string `json:"player_name"`
}
