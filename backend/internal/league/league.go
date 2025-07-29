package league

import (
	"context"
	"ecpc-league/engines"
	"ecpc-league/internal/messages"
	"encoding/json"
	"fmt"
	"sort"
)

func Create(ctx context.Context, leagueName, description string, scoring []messages.LeagueScoring) error {

	err := checkScoringIsSane(scoring)

	if err != nil {
		return err
	}

	tx := engines.MustTxFromContext(ctx)
	var leagueID int64

	err = tx.QueryRowContext(ctx, `
		INSERT INTO league.league (league_name, description) VALUES ($1, $2)	
		RETURNING id_league
	`, leagueName, description).Scan(&leagueID)

	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO league.league_scoring (id_league, rank, points)
		VALUES ($1, $2, $3)
	`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range scoring {
		_, err := stmt.ExecContext(ctx, leagueID, s.Rank, s.Score)
		if err != nil {
			return err
		}
	}

	return nil
}

func List(ctx context.Context) (messages.ListLeagueResponse, error) {

	tx := engines.MustTxFromContext(ctx)

	rows, err := tx.QueryContext(ctx, `
		SELECT id_league, league_name, description FROM league.league
	`)

	if err != nil {
		return messages.ListLeagueResponse{}, err
	}

	var res messages.ListLeagueResponse

	for rows.Next() {
		var l messages.ListLeagueResponseItem
		if err := rows.Scan(&l.IdLeague, &l.LeagueName, &l.Description); err != nil {
			return messages.ListLeagueResponse{}, err
		}
		res.Leagues = append(res.Leagues, l)
	}

	return res, nil
}

func Get(ctx context.Context, idLeague int64) (messages.GetLeagueResponse, error) {
	tx := engines.MustTxFromContext(ctx)

	row := tx.QueryRowContext(ctx, `
		SELECT
			l.id_league,
			l.league_name,
			l.description,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'rank', s.rank,
				'score', s.score
			)) FILTER (WHERE s.rank IS NOT NULL), '[]') AS scoring,

			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'id_league_round', r.id_league_round,
				'round_name', r.round_name
			)) FILTER (WHERE r.id_league_round IS NOT NULL), '[]') AS rounds,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'player_name', p.player_name,
				'id_league_player', p.id_league_player
			)) FILTER (WHERE p.player_name IS NOT NULL), '[]') AS players
		FROM league.league l
		LEFT JOIN league.league_scoring s USING(id_league)
		LEFT JOIN league.league_round r USING(id_league)
		LEFT JOIN league.league_player p USING(id_league)
		WHERE l.id_league = $1

		GROUP BY l.id_league, l.league_name, l.description;

	`, idLeague)

	var (
		id                                   int
		name, desc                           string
		scoringJSON, roundsJSON, playersJSON []byte
	)

	err := row.Scan(&id, &name, &desc, &scoringJSON, &roundsJSON, &playersJSON)

	if err != nil {
		return messages.GetLeagueResponse{}, err
	}

	var (
		scoring []messages.LeagueScoring
		rounds  []messages.LeagueRounds
		players []messages.LeaguePlayers
	)

	if err := json.Unmarshal(scoringJSON, &scoring); err != nil {
		return messages.GetLeagueResponse{}, fmt.Errorf("error unmarshalling scoring: %w", err)
	}
	if err := json.Unmarshal(roundsJSON, &rounds); err != nil {
		return messages.GetLeagueResponse{}, fmt.Errorf("error unmarshalling rounds: %w", err)
	}
	if err := json.Unmarshal(playersJSON, &players); err != nil {
		return messages.GetLeagueResponse{}, fmt.Errorf("error unmarshalling players: %w", err)
	}

	return messages.GetLeagueResponse{
		League: messages.GetLeagueItem{
			IdLeague:    id,
			LeagueName:  name,
			Description: desc,
			Scoring:     scoring,
			Rounds:      rounds,
			Players:     players,
		},
	}, nil

}

func getLeagueScoreboard(players []messages.LeaguePlayers, scoring []messages.LeagueScoring, results []messages.LeagueRoundsResults) []messages.LeagueScoreboard {
	playersScores := make(map[int]messages.LeagueScoreboard, 0)
	scoringMap := make(map[int]int, 0)

	for _, player := range players {
		playersScores[player.IdLeaguePlayer] = messages.LeagueScoreboard{
			Rank:         0,
			PlayerName:   player.PlayerName,
			RoundsPlayed: 0,
			TotalPoints:  0,
		}
	}

	for _, score := range scoring {
		scoringMap[score.Rank] = score.Score
	}

	for _, results := range results {
		cur := playersScores[results.IdLeaguePlayer]
		cur.RoundsPlayed += 1
		cur.TotalPoints += scoringMap[results.Rank]
	}

	var res []messages.LeagueScoreboard

	for _, value := range playersScores {
		res = append(res, value)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].TotalPoints > res[j].TotalPoints
	})

	if len(res) == 0 {
		return nil
	}

	res[0].Rank = 1

	for i := 1; i < len(res); i++ {
		if res[i].TotalPoints == res[i-1].TotalPoints {
			res[i].Rank = res[i-1].Rank
		} else {
			res[i].Rank = i + 1
		}
	}

	return res
}

func GetScoreboard(ctx context.Context, idLeague int64) (messages.GetLeagueScoreboardResponse, error) {

	tx := engines.MustTxFromContext(ctx)

	row := tx.QueryRowContext(ctx, `
		SELECT
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'rank', s.rank,
				'score', s.score
			)), '[]') AS scoring,

			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'id_league_round', r.id_league_round,
				'round_name', r.round_name
			)), '[]') AS rounds,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'player_name', p.player_name,
				'id_league_player', p.id_league_player
			)), '[]') AS players,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'id_league_round', rs.id_league_round,
				'id_league_player', rs.id_league_player,
				'rank', rs.rank
			)), '[]') as rounds_scores
		FROM league.league l
		LEFT JOIN league.league_scoring s USING(id_league)
		LEFT JOIN league.league_round r USING(id_league)
		LEFT JOIN league.league_player p USING(id_league)
		LEFT JOIN league.round_socre rs USING(id_league_round)
		WHERE l.id_league = $1

		GROUP BY l.id_league, l.league_name, l.description;

	`, idLeague)

	var scoringJSON, roundsJSON, playersJSON, roundsResultsJSON []byte

	var (
		scoring       []messages.LeagueScoring
		rounds        []messages.LeagueRounds
		players       []messages.LeaguePlayers
		roundsResults []messages.LeagueRoundsResults
	)

	err := row.Scan(&scoringJSON, &roundsJSON, &playersJSON, &roundsResultsJSON)

	if err != nil {
		return messages.GetLeagueScoreboardResponse{}, err
	}

	if err := json.Unmarshal(scoringJSON, &scoring); err != nil {
		return messages.GetLeagueScoreboardResponse{}, fmt.Errorf("error unmarshalling scoring: %w", err)
	}
	if err := json.Unmarshal(roundsJSON, &rounds); err != nil {
		return messages.GetLeagueScoreboardResponse{}, fmt.Errorf("error unmarshalling rounds: %w", err)
	}
	if err := json.Unmarshal(playersJSON, &players); err != nil {
		return messages.GetLeagueScoreboardResponse{}, fmt.Errorf("error unmarshalling players: %w", err)
	}
	if err := json.Unmarshal(roundsResultsJSON, &roundsResults); err != nil {
		return messages.GetLeagueScoreboardResponse{}, fmt.Errorf("error unmarshalling players: %w", err)
	}

	return messages.GetLeagueScoreboardResponse{
		RoundsResults: roundsResults,
		Scoreboard:    getLeagueScoreboard(players, scoring, roundsResults),
		Rounds:        rounds,
		Players:       players,
	}, nil

}
