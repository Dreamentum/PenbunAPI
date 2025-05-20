package utils

import (
	"database/sql"
)

func ExecuteTransaction(db *sql.DB, queries []func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	for _, query := range queries {
		if err := query(tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
