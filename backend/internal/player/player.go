package player

import (
	"context"
	"ecpc-league/engines"
	"ecpc-league/internal/league"
	"fmt"
)

func Create(ctx context.Context, idLeague int, playerName string) error {

	league, err := league.Get(ctx, int64(idLeague))

	if err != nil {
		return err
	}

	if league.League.IdLeague == 0 {
		return fmt.Errorf("this league doesn't exist")
	}

	tx := engines.MustTxFromContext(ctx)

	_, err = tx.ExecContext(ctx, `
		INSERT INTO league.player (id_league, player_name) VALUES ($1, $2)	
	`, idLeague, playerName)

	if err != nil {
		return err
	}

	return nil
}
