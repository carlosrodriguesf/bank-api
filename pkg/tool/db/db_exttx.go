package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type extendedTx struct {
	*sqlx.Tx
}

func NewExtendedTx(tx *sqlx.Tx) ExtendedTx {
	return &extendedTx{
		Tx: tx,
	}
}

func (db *extendedTx) NamedGetContext(ctx context.Context, query string, dest, arg interface{}) error {
	stmt, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = stmt.GetContext(ctx, dest, arg)
	if err != nil {
		return err
	}
	return nil
}

func (db *extendedTx) NamedSelectContext(ctx context.Context, query string, dest, arg interface{}) error {
	stmt, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = stmt.SelectContext(ctx, dest, arg)
	if err != nil {
		return err
	}
	return nil
}

func (db *extendedTx) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	stmt, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return stmt.QueryxContext(ctx, arg)
}
func (db *extendedTx) Close() error {
	return nil
}
