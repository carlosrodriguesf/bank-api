package db

import (
	"context"
	"errors"
)

func BeginTransaction(ctx context.Context, conn ExtendedDB) (ExtendedTx, error) {
	db, ok := conn.(*extendedDB)
	if !ok {
		return nil, errors.New("invalid connection")
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return NewExtendedTx(tx), nil
}
