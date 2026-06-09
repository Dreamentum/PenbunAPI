package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"PenbunAPI/config"
)

type TransactionStep struct {
	Name      string
	Query     string
	Args      []interface{}
	RowsAffected int64
}

func ExecuteTransaction(steps []TransactionStep) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	start := time.Now()
	config.TransactionLogger.Printf("TX START | steps=%d", len(steps))

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			config.TransactionLogger.Printf("TX ROLLBACK (panic) | duration=%v | panic=%v", time.Since(start), r)
			log.Printf("Transaction panic recovered: %v", r)
		}
	}()

	for i, step := range steps {
		stepStart := time.Now()

		result, err := tx.Exec(step.Query, step.Args...)
		if err != nil {
			tx.Rollback()
			config.TransactionLogger.Printf("TX ROLLBACK | step=%d/%d name=%s duration=%v error=%s",
				i+1, len(steps), step.Name, time.Since(stepStart), err)
			return fmt.Errorf("step %d (%s): %w", i+1, step.Name, err)
		}

		affected, _ := result.RowsAffected()
		steps[i].RowsAffected = affected

		config.TransactionLogger.Printf("TX STEP OK | step=%d/%d name=%s duration=%v rows=%d",
			i+1, len(steps), step.Name, time.Since(stepStart), affected)
	}

	if err := tx.Commit(); err != nil {
		config.TransactionLogger.Printf("TX COMMIT FAIL | duration=%v error=%s", time.Since(start), err)
		return fmt.Errorf("commit transaction: %w", err)
	}

	config.TransactionLogger.Printf("TX COMMIT OK | duration=%v steps=%d", time.Since(start), len(steps))
	return nil
}

func ScanRow(row *sql.Row, dest ...interface{}) error {
	if err := row.Scan(dest...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
