package engines

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
)

func GetConn(ctx *context.Context, dbUrl string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(*ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type txKey struct{}

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func TxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}

func MustTxFromContext(ctx context.Context) *sql.Tx {
	tx, ok := TxFromContext(ctx)
	if !ok {
		panic("Transaction doesn't exist")
	}
	return tx
}
