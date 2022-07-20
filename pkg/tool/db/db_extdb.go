package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type extendedDB struct {
	*sqlx.DB
}

func NewExtendedDB(db *sqlx.DB) ExtendedDB {
	return &extendedDB{
		DB: db,
	}
}

func (db *extendedDB) NamedGetContext(ctx context.Context, query string, dest, arg interface{}) error {
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

func (db *extendedDB) NamedSelectContext(ctx context.Context, query string, dest, arg interface{}) error {
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
