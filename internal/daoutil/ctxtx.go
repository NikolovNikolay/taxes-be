package daoutil

import (
	"context"
	"database/sql"
)

type ctxTxKey struct{}

func EnsureTransaction(ctx context.Context, db *sql.DB, retryableFn func(tx *sql.Tx) error) error {
	tx, ok := ctx.Value(ctxTxKey{}).(*sql.Tx)
	if !ok {
		newTx, err := db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})

		if err != nil {
			return err
		}
		tx = newTx
	}
	err := retryableFn(tx)
	if err != nil {
		return err
	}
	//err = tx.Commit()
	//if err != nil {
	//	err = tx.Rollback()
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

func InTransaction(ctx context.Context, db *sql.DB, retryableFn func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}
	tCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	tCtx = context.WithValue(tCtx, ctxTxKey{}, tx)
	err = retryableFn(tCtx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
	}

	return nil
}
