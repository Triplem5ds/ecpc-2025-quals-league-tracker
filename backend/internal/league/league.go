package league

import (
	"context"
	"ecpc-league/engines"
	"ecpc-league/internal/messages"
)

func Create(ctx context.Context, leagueName, url string) error {
	tx := engines.MustTxFromContext(ctx)

	_, err := tx.ExecContext(ctx, `
		INSERT INTO league.league (league_name, url) VALUES ($1, $2)	
	`, leagueName, url)

	if err != nil {
		return err
	}

	return nil
}

func List(ctx context.Context) (messages.ListLeagueResponse, error) {

	tx := engines.MustTxFromContext(ctx)

	rows, err := tx.QueryContext(ctx, `
		SELECT league_name, url FROM league.league
	`)

	if err != nil {
		return messages.ListLeagueResponse{}, err
	}

	var res messages.ListLeagueResponse

	for rows.Next() {
		var l messages.ListLeagueResponseItem
		if err := rows.Scan(&l.LeagueName, &l.URL); err != nil {
			return messages.ListLeagueResponse{}, err
		}
		res.Leagues = append(res.Leagues, l)
	}

	return res, nil
}
