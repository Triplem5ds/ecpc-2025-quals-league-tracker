package round

import (
	"context"
	"ecpc-league/engines"
	"ecpc-league/internal/league"
	"ecpc-league/internal/messages"
	"encoding/json"
	"fmt"
	"time"
)

func Create(ctx context.Context, idLeague int, roundName, url, description string, startTime time.Time) error {

	league, err := league.Get(ctx, int64(idLeague))

	if err != nil {
		return err
	}

	if league.League.IdLeague == 0 {
		return fmt.Errorf("this league doesn't exist")
	}

	tx := engines.MustTxFromContext(ctx)

	_, err = tx.ExecContext(ctx, `
		INSERT INTO league.round (id_league, round_name, start_time, url, description) VALUES ($1, $2, $3, $4, $5)	
	`, idLeague, roundName, startTime, url, description)

	if err != nil {
		return err
	}

	return nil
}

func Get(ctx context.Context, idLeagueRound int64) (messages.GetRoundResponse, error) {
	tx := engines.MustTxFromContext(ctx)

	row := tx.QueryRowContext(ctx, `
		SELECT
			r.id_leagueRound,
			r.round_name,
			r.description,
			r.start_time,
			r.url
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'id_league_player', s.id_league_player,
				'rank', s.player,
				'player_name', p.player_name
			)) FILTER (WHERE s.id_league_player IS NOT NULL), '[]') AS scoring
		FROM league.league_round r
		LEFT JOIN league.round_score s USING(id_league_round)
		LEFT JOIN league.league_player p USING(id_league_player)
		WHERE r.id_league_round = $1

		GROUP BY r.id_leagueRound, r.round_name, r.description, r.start_time, r.url;

	`, idLeagueRound)

	var (
		id              int
		name, desc, url string
		start_time      time.Time
		scoringJSON     []byte
	)

	err := row.Scan(&id, &name, desc, &start_time, &url, &scoringJSON)

	if err != nil {
		return messages.GetRoundResponse{}, err
	}

	var scoring []messages.RoundScoring

	if err := json.Unmarshal(scoringJSON, &scoring); err != nil {
		return messages.GetRoundResponse{}, fmt.Errorf("error unmarshalling scoring: %w", err)
	}

	return messages.GetRoundResponse{
		Round: messages.GetRoundItem{
			IdLeagueRound: id,
			RoundName:     name,
			Description:   desc,
			URL:           url,
			StartTime:     start_time,
			Scoring:       scoring,
		},
	}, nil

}

func Update(ctx context.Context, idLeagueRound int, scoring []messages.RoundScoring) error {
	round, err := Get(ctx, int64(idLeagueRound))

	if err != nil {
		return err
	}

	if round.Round.IdLeagueRound == 0 {
		return fmt.Errorf("this Round doesn't exist")
	}

	tx := engines.MustTxFromContext(ctx)

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO league.round_scoring (id_league_round, id_league_player, rank)
		VALUES ($1, $2, $3)
		ON CONFLICT (id_league_round, id_league_player)
		DO UPDATE SET rank = EXCLUDED.rank`)

	if err != nil {
		return err
	}

	for _, s := range scoring {
		_, err := stmt.ExecContext(ctx, idLeagueRound, s.IdLeaguePlayer, s.Rank)

		if err != nil {
			return err
		}

	}

	return nil
}
